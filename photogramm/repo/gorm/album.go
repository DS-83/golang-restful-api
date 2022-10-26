package gorm

import (
	"context"
	"errors"
	e "example-restful-api-server/err"
	"example-restful-api-server/models"

	"gorm.io/gorm"
)

type AlbumRepo struct {
	db *gorm.DB
}

type Album struct {
	ID        int
	AlbumName string
	UserID    int
	Photos    []Photo
}

func NewAlbumRepo(db *gorm.DB) *AlbumRepo {
	return &AlbumRepo{
		db: db,
	}
}

func (r *AlbumRepo) CreateAlbum(c context.Context, u *models.User, name string) error {
	a := &Album{AlbumName: name,
		UserID: u.Id,
	}
	// Check album name does not exist
	err := r.db.WithContext(c).Where("album_name = ? AND user_id = ?", a.AlbumName, a.UserID).First(a).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return e.ErrAlreadyExist
	}

	if err := r.db.Create(a).Error; err != nil {
		return err
	}

	return nil
}
func (r *AlbumRepo) GetAlbum(c context.Context, u *models.User, name string) (albm *models.PhotoAlbum, err error) {
	a := &Album{}
	if err = r.db.WithContext(c).Preload("Photos").Where("album_name = ? AND user_id = ?", name, u.Id).Find(a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrNotFound
		}
		return nil, err
	}
	albm = toModelsAlbum(a)

	return albm, nil
}

func (r *AlbumRepo) RemoveAlbum(c context.Context, u *models.User, name string) error {
	a := &Album{AlbumName: name,
		UserID: u.Id,
	}
	// Check existence
	if err := r.db.WithContext(c).Where("album_name = ? AND user_id = ?", name, u.Id).First(a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.ErrNotFound
		}
		return err
	}
	err := r.db.WithContext(c).Table("photos").Where("user_id = ?", a.UserID).Delete(a).Error
	if err != nil {
		return err
	}
	err = r.db.WithContext(c).Where("user_id = ?", a.UserID).Delete(a).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *AlbumRepo) GetInfo(c context.Context, u *models.User) (*models.User, error) {
	albums := []models.PhotoAlbum{}
	names := []string{}

	err := r.db.WithContext(c).Select("album_name").Table("albums").Where("user_id = ?", u.Id).Find(&names).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrNotFound
		}
		return nil, err
	}

	for _, name := range names {
		album, err := r.GetAlbum(c, u, name)
		if err != nil {
			return nil, err
		}
		albums = append(albums, *album)
	}

	u.PhotoAlbums = albums
	return u, nil
}

func toModelsAlbum(a *Album) *models.PhotoAlbum {
	albm := &models.PhotoAlbum{
		Name:   a.AlbumName,
		UserId: a.UserID,
		Photos: []string{},
		Total:  0,
	}

	photos := []string{}

	for _, p := range a.Photos {
		photos = append(photos, p.ID)
	}

	albm.Photos = photos
	albm.Total = len(photos)

	return albm
}
