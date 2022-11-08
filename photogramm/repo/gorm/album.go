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

type album struct {
	ID        int
	AlbumName string
	UserID    int
	Photos    []photo
}

type user struct {
	ID       int
	MongoID  string
	Username string
}

func NewAlbumRepo(db *gorm.DB) *AlbumRepo {
	return &AlbumRepo{
		db: db,
	}
}

func (r *AlbumRepo) CreateAlbum(c context.Context, u *models.User, name string) error {
	user := toDbUser(u)

	a := &album{
		AlbumName: name,
		UserID:    user.ID,
	}
	// Check album name does not exist
	err := r.db.WithContext(c).Where("album_name = ? AND user_id = ?", a.AlbumName, user.ID).First(a).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return e.ErrAlreadyExist
	}

	if err := r.db.Create(a).Error; err != nil {
		return err
	}

	return nil
}
func (r *AlbumRepo) GetAlbum(c context.Context, u *models.User, name string) (albm *models.PhotoAlbum, err error) {
	a := &album{
		AlbumName: name,
	}

	user := toDbUser(u)
	if err := r.db.WithContext(c).Where(user).First(user).Error; err != nil {
		return nil, err
	}

	if err = r.db.WithContext(c).Preload("Photos").Where(
		"album_name = ? AND user_id = ?", a.AlbumName, user.ID).Take(a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrNotFound
		}
		return nil, err
	}

	return toModelsAlbum(a), nil
}

func (r *AlbumRepo) RemoveAlbum(c context.Context, u *models.User, name string) error {
	a := &album{
		AlbumName: name,
	}

	user := toDbUser(u)

	if err := r.db.WithContext(c).Where(user).First(user).Error; err != nil {
		return err
	}

	// Check existence
	if err := r.db.WithContext(c).Where("album_name = ? AND user_id = ?", name, user.ID).First(a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.ErrNotFound
		}
		return err
	}
	err := r.db.WithContext(c).Where("user_id = ? AND album_id = ?", u.ID, a.ID).Delete(&photo{}).Error
	if err != nil {
		return err
	}
	err = r.db.WithContext(c).Where("user_id = ?", user.ID).Delete(a).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *AlbumRepo) GetInfo(c context.Context, u *models.User) (*models.User, error) {
	albums := []models.PhotoAlbum{}
	names := []string{}

	user := toDbUser(u)

	if err := r.db.WithContext(c).Where(user).First(user).Error; err != nil {
		return nil, err
	}

	err := r.db.WithContext(c).Select("album_name").Table("albums").Where("user_id = ?", user.ID).Find(&names).Error
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

func toModelsAlbum(a *album) *models.PhotoAlbum {

	albm := &models.PhotoAlbum{
		Name:     a.AlbumName,
		UserID:   a.UserID,
		PhotosID: []string{},
		Total:    0,
	}

	photos := []string{}

	for _, p := range a.Photos {
		photos = append(photos, p.ID)
	}

	albm.PhotosID = photos
	albm.Total = len(photos)

	return albm
}

func toDbUser(u *models.User) *user {
	return &user{
		ID:       u.ID,
		MongoID:  u.MongoID,
		Username: u.Username,
	}
}
