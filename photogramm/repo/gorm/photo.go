package gorm

import (
	"context"
	"errors"
	"example-restful-api-server/e"
	"example-restful-api-server/models"
	"io"
	"log"

	"gorm.io/gorm"
)

type PhotoRepo struct {
	db           *gorm.DB
	defaultAlbum string
}

type User struct {
	ID       int
	Username string
	Password string
}

type Photo struct {
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
	err := r.db.WithContext(ctx).Select("albums.id").Joins("INNER JOIN users ON albums.user_id=users.id").Where("username= ? AND album_name = ?",
		p.Username, p.AlbumName).Model(&Album{}).First(&photo.AlbumID).Error
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
	photo := &Photo{ID: id}

	err = r.db.WithContext(ctx).Take(photo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, e.ErrNotFound
		}
		return nil, err
	}

	p.Id = photo.ID

	if err = r.db.WithContext(ctx).Model(&Album{}).Select(
		"AlbumName", "id = ?", photo.AlbumID).Take(&p.AlbumName).Error; err != nil {
		return nil, err
	}

	if err = r.db.WithContext(ctx).Model(&User{}).Select(
		"Username", "id = ?", photo.UserID).Take(&p.Username).Error; err != nil {
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

func toGormPhoto(p *models.Photo) *Photo {
	return &Photo{
		ID:     p.Id,
		UserID: p.UserId,
	}
}
