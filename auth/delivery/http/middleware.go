package http

import (
	"example-restful-api-server/auth"
	"example-restful-api-server/e"
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
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user, err := m.uc.ParseTokenFromString(tokenString)
	if err != nil {
		status := http.StatusInternalServerError
		if err == e.ErrInvalidAccessToken {
			status = http.StatusUnauthorized
		}

		log.Println(err)
		c.AbortWithStatus(status)
		return
	}
	c.Set(auth.CtxUserKey, user)
}

func parseHeaderAuth(h http.Header) (key string, err error) {
	defer func() { err = e.Wrap("can't read authorization param", err) }()
	header := h["Authorization"]
	if len(header) == 0 {
		return key, fmt.Errorf("missing params")
	}

	s := strings.Fields(header[0])
	if len(s) < 2 || len(s) > 2 {
		return key, fmt.Errorf("incorrect")
	}

	if s[0] != authBearer {
		return key, fmt.Errorf("incorrect auth type")
	}
	return s[1], nil

}
