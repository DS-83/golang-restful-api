package gorm

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"example-restful-api-server/e"
	"example-restful-api-server/models"
	"log"

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
		Id:          u.ID,
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

type token struct {
	ID string `gorm:"primaryKey"`
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
	err := r.db.WithContext(ctx).Where("user_id = ?", u.Id).Delete(&photo{}).Error
	if err != nil {
		return err
	}
	err = r.db.WithContext(ctx).Where("user_id = ?", u.Id).Delete(&album{}).Error
	if err != nil {
		return err
	}

	err = r.db.WithContext(ctx).Delete(&user{ID: u.Id}).Error

	return err
}

func (r UserRepo) RevokeToken(ctx context.Context, key []byte) error {
	// Create hash of key string
	hasher := sha1.New()
	if _, err := hasher.Write(key); err != nil {
		log.Println(err)
		return err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	err := r.db.WithContext(ctx).Table("revoked_tokens").First(&token{ID: hash}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err)
		return e.ErrRevokedToken
	}

	if err = r.db.WithContext(ctx).Table("revoked_tokens").Create(&token{ID: hash}).Error; err != nil {
		return err
	}

	return nil
}

func (r UserRepo) IsRevoked(key []byte) bool {
	// Create hash of key string
	hasher := sha1.New()
	if _, err := hasher.Write(key); err != nil {
		log.Println(err)
		return true
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	err := r.db.Table("revoked_tokens").First(&token{ID: hash}).Error

	return !errors.Is(err, gorm.ErrRecordNotFound)
}
