// nolint
package mock

import (
	"context"
	"example-restful-api-server/models"

	"github.com/stretchr/testify/mock"
)

type UserStorageMock struct {
	mock.Mock
}

func (s *UserStorageMock) CreateUser(ctx context.Context, user *models.User) error {
	args := s.Called(user)

	return args.Error(0)
}

func (s *UserStorageMock) GetUser(ctx context.Context, username, password string) (*models.User, error) {
	args := s.Called(username, password)

	return args.Get(0).(*models.User), args.Error(1)
}

func (s *UserStorageMock) DeleteUser(c context.Context, u *models.User) error {
	args := s.Called(u)

	return args.Error(0)
}

func (s *UserStorageMock) UpdateUser(c context.Context, f, u *models.User) error {
	args := s.Called(f, u)

	return args.Error(0)
}
