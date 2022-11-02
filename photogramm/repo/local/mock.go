package local

import (
	"context"
	"example-restful-api-server/models"
	"io"

	"github.com/stretchr/testify/mock"
)

type PhotoLocalStorageMock struct {
	mock.Mock
}

func (s *PhotoLocalStorageMock) CreatePhoto(c context.Context, p *models.Photo, r io.Reader) (string, error) {
	args := s.Called(p)
	return args.String(0), args.Error(1)
}
func (s *PhotoLocalStorageMock) GetPhoto(c context.Context, u *models.User, id string) (*models.Photo, error) {
	args := s.Called(u, id)
	return args.Get(0).(*models.Photo), args.Error(1)
}
func (s *PhotoLocalStorageMock) RemovePhoto(c context.Context, p *models.Photo) error {
	args := s.Called(p)
	return args.Error(0)
}
