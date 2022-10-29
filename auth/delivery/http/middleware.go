package http

import (
	"example-restful-api-server/auth"
	e "example-restful-api-server/err"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Auth types
const (
	authBasic  string = "Basic"
	authBearer string = "Bearer"
	authHeader string = "Authorization"
)

type AuthMiddleware struct {
	uc auth.UseCase
}

func NewAuthMiddleware(usecase auth.UseCase) gin.HandlerFunc {
	return (&AuthMiddleware{
		uc: usecase,
	}).Handle
}

// Middleware get jwt token from request header, parse and set user
func (m *AuthMiddleware) Handle(c *gin.Context) {

	tokenString, err := parseHeaderAuth(c.Request.Header)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := m.uc.ParseTokenFromString(c, &tokenString)
	if err != nil {
		status := http.StatusInternalServerError
		if err == e.ErrInvalidAccessToken {
			status = http.StatusUnauthorized
		}

		log.Println(err)
		c.AbortWithStatus(status)
		return
	}
	c.Set(auth.CtxTokenString, tokenString)
	c.Set(auth.CtxUserKey, user)
}

func parseHeaderAuth(h http.Header) (key string, err error) {
	defer func() { err = e.Wrap("can't read authorization param", err) }()
	header := h[authHeader]
	if len(header) == 0 {
		return key, fmt.Errorf("missing params")
	}

	s := strings.Fields(header[0])
	if len(s) != 2 {
		return key, fmt.Errorf("incorrect auth string")
	}

	if s[0] != authBearer {
		return key, fmt.Errorf("incorrect auth type")
	}
	return s[1], nil

}
