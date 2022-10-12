package auth

import (
	"context"
	"fotogramm/example-restful-api-server/models"
)

// Users storage interface
type UserRepo interface {
	CreateUser(context.Context, *models.User) error
	GetUser(c context.Context, u string, p string) (*models.User, error)
	// UpdateUser(context.Context, *model.DBUser) (*model.DBUser, error)
	DeleteUser(context.Context, *models.User) error
	RevokeToken(ctx context.Context, key []byte) error
	IsRevoked(key []byte) bool
}
