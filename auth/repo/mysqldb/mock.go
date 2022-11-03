// nolint
package mysqldb

import (
	"context"
	"example-restful-api-server/models"

	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m UserRepoMock) CreateUser(ctx context.Context, u *models.User) error {

	args := m.Called(u)

	return args.Error(0)
}

func (m UserRepoMock) GetUser(ctx context.Context, username string, pass string) (*models.User, error) {
	args := m.Called(username, pass)

	return args.Get(0).(*models.User), args.Error(1)
}

func (m UserRepoMock) DeleteUser(ctx context.Context, u *models.User) error {
	args := m.Called(u)

	return args.Error(0)
}

func (m UserRepoMock) IsRevoked(key []byte) bool {
	args := m.Called(key)

	return args.Bool(0)
}

func (m UserRepoMock) RevokeToken(ctx context.Context, key []byte) error {
	args := m.Called(key)

	return args.Error(0)
}
