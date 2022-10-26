package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type TokenStorageMock struct {
	mock.Mock
}

func (s *TokenStorageMock) RevokeToken(c context.Context, t *string) error {
	args := s.Called(*t)
	return args.Error(0)
}

func (s *TokenStorageMock) IsRevoked(c context.Context, t *string) (bool, error) {
	args := s.Called(*t)

	return args.Bool(0), args.Error(1)
}
