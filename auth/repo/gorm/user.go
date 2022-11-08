package gorm

import (
	"context"
	"example-restful-api-server/models"

	"gorm.io/gorm"
)

type user struct {
	ID       int
	MongoID  string
	Username string
}

func NewUserRepo(db *gorm.DB, defaultAlbum string) *UserRepo {
	return &UserRepo{
		db:           db,
		defaultAlbum: defaultAlbum,
	}
}

type UserRepo struct {
	db           *gorm.DB
	defaultAlbum string
}

type album struct {
	ID        int
	AlbumName string
	UserID    int
}

type photo struct {
	ID      string `gorm:"primaryKey"`
	UserID  int
	AlbumID int
}

func (r UserRepo) CreateUser(ctx context.Context, u *models.User) (*models.User, error) {
	user := toGormUser(u)
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	// Create default album for new user
	album := album{
		UserID:    user.ID,
		AlbumName: r.defaultAlbum,
	}
	if err := r.db.WithContext(ctx).Table("albums").Create(&album).Error; err != nil {
		return nil, err
	}

	return toModelUser(user), nil
}

func (r UserRepo) GetUser(ctx context.Context, username string) (*models.User, error) {
	user := user{
		Username: username,
	}
	if err := r.db.WithContext(ctx).Where(
		"username = ?", user.Username).First(&user).Error; err != nil {
		return nil, err
	}

	return toModelUser(&user), nil
}

func (r UserRepo) DeleteUser(ctx context.Context, u *models.User) error {
	user := toGormUser(u)

	if err := r.db.WithContext(ctx).Where(user).First(user).Error; err != nil {
		return err
	}

	err := r.db.WithContext(ctx).Where("user_id = ?", user.ID).Delete(&photo{}).Error
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).Where("user_id = ?", user.ID).Delete(&album{}).Error
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).Delete(user).Error

	return err
}

func (r UserRepo) UpdateUser(c context.Context, f, u *models.User) error {
	filt := toGormUser(f)
	upd := toGormUser(u)

	if err := r.db.WithContext(c).Where(filt).First(filt).Error; err != nil {
		return err
	}

	if err := r.db.Model(filt).Updates(upd).Error; err != nil {
		return err
	}
	return nil
}

func toGormUser(u *models.User) *user {
	return &user{
		ID:       u.ID,
		MongoID:  u.MongoID,
		Username: u.Username,
	}
}
func toModelUser(u *user) *models.User {
	return &models.User{
		ID:          u.ID,
		MongoID:     u.MongoID,
		Username:    u.Username,
		PhotoAlbums: []models.PhotoAlbum{},
	}
}
