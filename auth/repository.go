package auth

import (
	"context"
	"example-restful-api-server/models"
)

// Users storage interface
type UserRepo interface {
	CreateUser(context.Context, *models.User) (*models.User, error)
	GetUser(c context.Context, u string) (*models.User, error)
	UpdateUser(c context.Context, f, u *models.User) error
	DeleteUser(context.Context, *models.User) error
}
