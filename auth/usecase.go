package auth

import (
	"context"
	"example-restful-api-server/models"
)

const CtxUserKey = "user"

type UseCase interface {
	SignUp(c context.Context, u string, p string) error
	SignIn(c context.Context, u string, p string) (*string, error)
	DeleteUser(c context.Context, u *models.User) error
	ParseTokenFromString(tokenString string) (*models.User, error)
}
