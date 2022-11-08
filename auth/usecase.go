package auth

import (
	"context"
	"example-restful-api-server/models"
)

const (
	CtxUserKey     = "user"
	CtxTokenString = "token"
)

type UseCase interface {
	SignUp(c context.Context, u string, p string) error
	SignIn(c context.Context, u string, p string) (*string, error)
	DeleteUser(c context.Context, u *models.User, t string) error
	ParseTokenFromString(c context.Context, tokenString string) (*models.User, error)
}
