package gorm

import (
	"context"
	"errors"
	e "example-restful-api-server/err"
	"example-restful-api-server/models"
	"io"
	"log"

	"gorm.io/gorm"
)

type PhotoRepo struct {
	db           *gorm.DB
	defaultAlbum string
}

type photo struct {
	ID      string `gorm:"primaryKey"`
	UserID  int
	AlbumID int
}

func NewPhotoRepo(db *gorm.DB, defaultAlbum string) *PhotoRepo {
	return &PhotoRepo{
		db:           db,
		defaultAlbum: defaultAlbum,
	}
}

func (r *PhotoRepo) CreatePhoto(ctx context.Context, p *models.Photo, s io.Reader) (string, error) {
	if p.AlbumName == "" {
		p.AlbumName = r.defaultAlbum
	}

	photo := toGormPhoto(p)
	err := r.db.WithContext(ctx).Select("albums.id").Where("user_id= ? AND album_name = ?",
		p.UserID, p.AlbumName).Model(&album{}).First(&photo.AlbumID).Error
	if err != nil {
		log.Println(err)
		return "", err
	}
	if err = r.db.WithContext(ctx).Create(photo).Error; err != nil {
		return "", err
	}
	return photo.ID, nil
}

func (r *PhotoRepo) GetPhoto(ctx context.Context, u *models.User, id string) (p *models.Photo, err error) {
	photo := &photo{ID: id}

	err = r.db.WithContext(ctx).Take(photo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrNotFound
		}
		return nil, err
	}

	p = new(models.Photo)
	p.ID = id
	if err = r.db.WithContext(ctx).Model(&album{}).Select(
		"album_name").Where("id=?", photo.AlbumID).Take(&p.AlbumName).Error; err != nil {
		return nil, err
	}

	if err = r.db.WithContext(ctx).Model(&user{}).Select(
		"username").Where("id=?", photo.UserID).Take(&p.Username).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PhotoRepo) RemovePhoto(ctx context.Context, p *models.Photo) error {
	photo := toGormPhoto(p)

	err := r.db.WithContext(ctx).Where("user_id = ?", photo.UserID).Delete(&photo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return e.ErrNotFound
		}
		return err
	}

	return nil
}

func toGormPhoto(p *models.Photo) *photo {
	return &photo{
		ID:     p.ID,
		UserID: p.UserID,
	}
}
