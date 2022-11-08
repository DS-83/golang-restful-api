package auth

import (
	"context"
	"example-restful-api-server/models"
)

// Users storage interface
type UserRepo interface {
	CreateUser(context.Context, *models.User) error
	GetUser(c context.Context, u string, p string) (*models.User, error)
	// UpdateUser(context.Context, *model.DBUser) (*model.DBUser, error)
	DeleteUser(context.Context, *models.User) error
	// RevokeToken(c context.Context, key []byte) error
	// IsRevoked(key []byte) bool
}

// Tokens storage interface
type TokenRepo interface {
	RevokeToken(c context.Context, t string) error
	IsRevoked(c context.Context, t string) (bool, error)
}
