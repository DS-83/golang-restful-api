package usecase

import (
	"context"
	"example-restful-api-server/auth"
	"example-restful-api-server/e"
	"example-restful-api-server/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	userRepo auth.UserRepo
	jwtKey   []byte
}

type AuthClaims struct {
	User *models.User
	jwt.RegisteredClaims
}

func NewAuthUsecase(a auth.UserRepo, b []byte) *AuthUsecase {
	return &AuthUsecase{
		userRepo: a,
		jwtKey:   b,
	}
}

func (c *AuthUsecase) SignUp(ctx context.Context, username, pass string) (err error) {
	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8
	//  (this value can be more or less, depending on the computing power you wish to utilize)
	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(pass), 8)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
	}
	return c.userRepo.CreateUser(ctx, user)
}

// Sign in user and get JWT string
func (c *AuthUsecase) SignIn(ctx context.Context, username, pass string) (*string, error) {
	user, err := c.userRepo.GetUser(ctx, username, pass)
	if err != nil {
		return nil, err
	}
	// Token expiration time:
	exp := time.Now().Add(86400 * time.Second)
	// Create the Claims
	claims := AuthClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	// Create jwt Token. Signing method HS256 uses a []byte key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Token string
	ts, err := token.SignedString(c.jwtKey)
	if err != nil {
		return nil, err
	}

	return &ts, nil
}

func (c *AuthUsecase) ParseTokenFromString(tokenString string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return c.jwtKey, nil
	})

	if err != nil {
		return nil, e.Wrap("ParseTokenFromString", err)
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok || !token.Valid {
		return nil, e.ErrInvalidAccessToken
	}
	// Check for revoked token
	if ok := c.userRepo.IsRevoked([]byte(tokenString)); ok {
		return nil, e.ErrInvalidAccessToken
	}

	return claims.User, nil
}

func (c *AuthUsecase) DeleteUser(ctx context.Context, u *models.User) error {
	if err := c.userRepo.DeleteUser(ctx, u); err != nil {
		return err
	}
	if err := c.userRepo.RevokeToken(ctx, c.jwtKey); err != nil {
		return err
	}
	return nil
}
