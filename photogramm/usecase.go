package photogramm

import (
	"context"
	"fotogramm/example-restful-api-server/models"
	"io"
)

type UseCase interface {
	UploadPhoto(context.Context, *models.User, string, io.Reader) error
	GetPhoto(c context.Context, u *models.User, id string) (*models.Photo, error)
	RemovePhoto(c context.Context, u *models.User, id string) error
	CreateAlbum(c context.Context, u *models.User, name string) error
	GetAlbum(c context.Context, u *models.User, name string) (*models.PhotoAlbum, error)
	RemoveAlbum(c context.Context, u *models.User, name string) error
	GetInfo(c context.Context, u *models.User) (*models.User, error)
}
