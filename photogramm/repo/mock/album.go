// nolint
package mock

import (
	"context"
	"example-restful-api-server/models"

	"github.com/stretchr/testify/mock"
)

type AlbumRepoMock struct {
	mock.Mock
}

func (r *AlbumRepoMock) CreateAlbum(c context.Context, u *models.User, name string) error {
	args := r.Called(u, name)
	return args.Error(0)
}

func (r *AlbumRepoMock) GetAlbum(c context.Context, u *models.User, name string) (*models.PhotoAlbum, error) {
	args := r.Called(u, name)
	return args.Get(0).(*models.PhotoAlbum), args.Error(1)
}
func (r *AlbumRepoMock) RemoveAlbum(c context.Context, u *models.User, name string) error {
	args := r.Called(u, name)
	return args.Error(0)
}
func (r *AlbumRepoMock) GetInfo(c context.Context, u *models.User) (*models.User, error) {
	args := r.Called(u)
	return args.Get(0).(*models.User), args.Error(1)
}
