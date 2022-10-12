package photogramm

import (
	"context"
	"fotogramm/example-restful-api-server/models"
	"io"
)

// Photos storage inteface
type PhotoRepo interface {
	CreatePhoto(context.Context, *models.Photo, io.Reader) error
	GetPhoto(c context.Context, u *models.User, id string) (*models.Photo, error)
	RemovePhoto(context.Context, *models.Photo) error
}

// Photo albums storage interface
type PhotoAlbums interface {
	CreateAlbum(c context.Context, u *models.User, name string) error
	GetAlbum(c context.Context, u *models.User, name string) (*models.PhotoAlbum, error)
	RemoveAlbum(c context.Context, u *models.User, name string) error
	GetInfo(c context.Context, u *models.User) (*models.User, error)
}
