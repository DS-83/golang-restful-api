// nolint
package usecase

import (
	"context"
	"example-restful-api-server/models"
	"io"

	"github.com/stretchr/testify/mock"
)

type PhotogrammUsecaseMock struct {
	mock.Mock
}

func (m *PhotogrammUsecaseMock) UploadPhoto(c context.Context, p *models.Photo, i io.Reader) (string, error) {
	args := m.Called(p)

	return args.String(0), args.Error(1)
}

func (m *PhotogrammUsecaseMock) GetPhoto(c context.Context, u *models.User, id string) (*models.Photo, error) {
	args := m.Called(u, id)

	return args.Get(0).(*models.Photo), args.Error(1)
}

func (m *PhotogrammUsecaseMock) RemovePhoto(c context.Context, u *models.User, id string) error {
	args := m.Called(u, id)

	return args.Error(0)
}

func (m *PhotogrammUsecaseMock) CreateAlbum(c context.Context, u *models.User, name string) error {
	args := m.Called(u, name)

	return args.Error(0)
}

func (m *PhotogrammUsecaseMock) GetAlbum(c context.Context, u *models.User, name string) (*models.PhotoAlbum, error) {
	args := m.Called(u, name)

	return args.Get(0).(*models.PhotoAlbum), args.Error(1)
}

func (m *PhotogrammUsecaseMock) RemoveAlbum(c context.Context, u *models.User, name string) error {
	args := m.Called(u, name)

	return args.Error(0)
}

func (m *PhotogrammUsecaseMock) GetInfo(c context.Context, u *models.User) (*models.User, error) {
	args := m.Called(u)

	return args.Get(0).(*models.User), args.Error(1)
}
