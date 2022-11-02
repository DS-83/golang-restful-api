package usecase

import (
	"context"
	"example-restful-api-server/models"

	"github.com/stretchr/testify/mock"
)

type AuthUsecaseMock struct {
	mock.Mock
}

func (m *AuthUsecaseMock) SignUp(ctx context.Context, username, pass string) (err error) {
	args := m.Called(username, pass)

	return args.Error(0)
}

func (m *AuthUsecaseMock) SignIn(c context.Context, u string, p string) (*string, error) {
	args := m.Called(u, p)
	return args.Get(0).(*string), args.Error(1)
}

func (m *AuthUsecaseMock) DeleteUser(c context.Context, u *models.User, t string) error {
	args := m.Called(u, t)

	return args.Error(0)
}

func (m *AuthUsecaseMock) ParseTokenFromString(c context.Context, tokenString string) (*models.User, error) {
	args := m.Called(tokenString)

	return args.Get(0).(*models.User), args.Error(1)
}
