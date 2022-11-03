// nolint
package mock

import (
	"context"
	"example-restful-api-server/models"
	"io"

	"github.com/stretchr/testify/mock"
)

type PhotoRepoMock struct {
	mock.Mock
}

func (r *PhotoRepoMock) CreatePhoto(c context.Context, p *models.Photo, i io.Reader) (string, error) {
	args := r.Called(p)
	return args.String(0), args.Error(1)
}
func (r *PhotoRepoMock) GetPhoto(c context.Context, u *models.User, id string) (*models.Photo, error) {
	args := r.Called(u, id)
	return args.Get(0).(*models.Photo), args.Error(1)
}
func (r *PhotoRepoMock) RemovePhoto(c context.Context, p *models.Photo) error {
	args := r.Called(p)
	return args.Error(0)
}
