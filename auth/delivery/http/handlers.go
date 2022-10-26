package http

import (
	"encoding/base64"
	"example-restful-api-server/auth"
	e "example-restful-api-server/err"
	"example-restful-api-server/models"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase auth.UseCase
}

type userInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type deleteInput struct {
	Delete string `json:"delete"`
}

type response struct {
	Response string `json:"responce"`
}

type signInResp struct {
	Token string `json:"token"`
}

func NewHandler(uc auth.UseCase) *Handler {
	return &Handler{
		useCase: uc,
	}
}

func (h *Handler) SignUp(c *gin.Context) {
	input := userInput{}

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("sign-up: %s", err)
		return
	}
	if input.Username == "" || input.Password == "" {
		c.JSON(http.StatusNotAcceptable, response{Response: "incorrect input"})
		return
	}

	if err := h.useCase.SignUp(c.Request.Context(), input.Username, input.Password); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Printf("sign-up: %s", err)
		return
	}

	c.JSON(http.StatusOK, response{Response: "success"})

}

func (h *Handler) SignIn(c *gin.Context) {

	cred, err := parseSignInHeader(c.Request.Header)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("sign-in: %s", e.Wrap("can't read authorization param", err))
		return
	}

	if len(cred[0]) == 0 || len(cred[1]) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("sign-in: %s", fmt.Errorf("incorrect input"))
		return
	}
	token, err := h.useCase.SignIn(c.Request.Context(), cred[0], cred[1])
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		log.Printf("sign-in: %s", err)
		return
	}

	c.JSON(http.StatusOK, signInResp{Token: *token})
}

func (h *Handler) Delete(c *gin.Context) {

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	input := deleteInput{}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response{Response: "incorrect request body"})
		log.Printf("delete: %s", err)
		return
	}
	t := c.GetString(auth.CtxTokenString)
	err := h.useCase.DeleteUser(c, user, &t)
	if err == e.ErrRevokedToken {
		log.Printf("delete: %s", err)
		c.JSON(http.StatusUnauthorized, response{Response: "not valid token"})
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("delete: %s", err)
		return
	}

	c.JSON(http.StatusOK, response{Response: "delete success"})
}

func parseSignInHeader(h http.Header) (cred []string, err error) {
	header := h[authHeader]
	if len(header) == 0 {
		return nil, fmt.Errorf("missing params")
	}
	s := strings.Fields(header[0])
	if len(s) != 2 {
		return nil, fmt.Errorf("missing fields")
	}

	if s[0] != authBasic {
		return nil, e.ErrInvalidAuthType
	}

	data, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return nil, err
	}
	cred = strings.Split(string(data), ":")
	if len(cred) == 0 || len(cred) > 2 {
		return nil, fmt.Errorf("incorrect auth input")
	}

	return cred, nil

}
