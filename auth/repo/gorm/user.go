package gorm

import (
	"context"
	e "example-restful-api-server/err"
	"example-restful-api-server/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type user struct {
	ID       int
	Username string
	Password string
}

func NewUserRepo(db *gorm.DB, defaultAlbum string) *UserRepo {
	return &UserRepo{
		db:           db,
		defaultAlbum: defaultAlbum,
	}
}

func toGormUser(u *models.User) *user {
	return &user{
		Username: u.Username,
		Password: u.Password,
	}
}
func toModelUser(u *user) *models.User {
	return &models.User{
		ID:          u.ID,
		Username:    u.Username,
		PhotoAlbums: []models.PhotoAlbum{},
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

func (r UserRepo) CreateUser(ctx context.Context, u *models.User) (err error) {
	user := toGormUser(u)
	if err = r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}
	// Create default album for new user
	album := album{
		UserID:    user.ID,
		AlbumName: r.defaultAlbum,
	}
	err = r.db.WithContext(ctx).Table("albums").Create(&album).Error

	return err
}

func (r UserRepo) GetUser(ctx context.Context, username string, pass string) (*models.User, error) {
	user := user{
		Username: username,
	}
	if err := r.db.WithContext(ctx).Where(
		"username = ?", user.Username).First(&user).Error; err != nil {
		return nil, err
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
	if err != nil {
		return nil, e.Wrap("invalid credentials", err)
	}
	return toModelUser(&user), nil
}

func (r UserRepo) DeleteUser(ctx context.Context, u *models.User) error {
	err := r.db.WithContext(ctx).Where("user_id = ?", u.ID).Delete(&photo{}).Error
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).Where("user_id = ?", u.ID).Delete(&album{}).Error
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).Delete(&user{ID: u.ID}).Error

	return err
}
