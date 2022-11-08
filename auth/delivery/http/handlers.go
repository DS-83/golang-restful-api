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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

type updateInput struct {
	Update userInput `json:"update"`
}

type response struct {
	Response string `json:"response"`
}

type signInResp struct {
	Token string `json:"token"`
}

func NewHandler(uc auth.UseCase) *Handler {
	return &Handler{
		useCase: uc,
	}
}

// Sign up
// @Summary      Sign up user
// @Description  Register user based on login and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body userInput true "Username/Password"
// @Success      200  {object}  response
// @Failure      400
// @Failure      406  {object}  response
// @Failure      409  {object}  response
// @Failure      404
// @Failure      500
// @Router       /auth/sign-up [post]
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
		if status.Code(err) == codes.AlreadyExists {
			c.JSON(http.StatusConflict, response{Response: status.Convert(err).Message()})
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Printf("sign-up: %s", err)
		return
	}

	c.JSON(http.StatusOK, response{Response: "success"})

}

// Sign in
// @Summary      Sign in user
// @Description  Sign in user based on login and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security	 BasicAuth
// @Success      200  {object}  signInResp
// @Failure      400 {object} response
// @Failure      401 {object} response
// @Router       /auth/sign-in [post]
func (h *Handler) SignIn(c *gin.Context) {

	cred, err := parseSignInHeader(c.Request.Header)
	if err != nil {
		c.JSON(http.StatusBadRequest, response{Response: "incorrect request body"})
		log.Printf("sign-in: %s", e.Wrap("can't read authorization param", err))
		return
	}

	if len(cred[0]) == 0 || len(cred[1]) == 0 {
		c.JSON(http.StatusBadRequest, response{Response: "incorrect input"})
		log.Printf("incorrect input")
		return
	}
	token, err := h.useCase.SignIn(c.Request.Context(), cred[0], cred[1])
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			response{Response: status.Convert(err).Message()},
		)
		log.Printf("sign-in: %s", err)
		return
	}

	c.JSON(http.StatusOK, signInResp{Token: *token})
}

// Delete
// @Summary      Delete user
// @Description  Delete authorized user account
// @Tags         user
// @Accept       json
// @Produce      json
// @Security	 JWT
// @Param 		 delete body deleteInput true "delete input"
// @Success      200  {object}  response
// @Failure      400 {object} response
// @Failure      401 {object} response
// @Failure      500
// @Router       /user/delete [delete]
func (h *Handler) Delete(c *gin.Context) {

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	input := new(deleteInput)

	if err := c.BindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, response{Response: "incorrect request body"})
		log.Printf("delete: %s", err)
		return
	}
	t := c.GetString(auth.CtxTokenString)
	if err := h.useCase.DeleteUser(c, user, t); err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(
				http.StatusBadRequest,
				response{Response: status.Convert(err).Message()})
			log.Printf("delete: %s", err)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		log.Printf("delete: %s", err)
		return
	}

	c.JSON(http.StatusOK, response{Response: "delete success"})
}

// Update user
// @Summary      Update user
// @Description  Update username or password
// @Tags         user
// @Accept       json
// @Produce      json
// @Security	 JWT
// @Param 		 update body updateInput true "update"
// @Success      200  {object}  response
// @Failure      400 {object} response
// @Failure      401 {object} response
// @Failure      500
// @Router       /user/update [post]
func (h *Handler) Update(c *gin.Context) {

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	input := new(updateInput)

	if err := c.BindJSON(input); err != nil {
		c.JSON(http.StatusBadRequest, response{Response: "incorrect request body"})
		log.Println(err)
		return
	}
	upd := &models.User{
		Username: input.Update.Username,
		Password: input.Update.Password,
	}

	t := c.GetString(auth.CtxTokenString)

	if err := h.useCase.UpdateUser(c, user, upd, t); err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(
				http.StatusBadRequest,
				response{Response: status.Convert(err).Message()})
			log.Printf("delete: %s", err)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	c.JSON(http.StatusOK, response{Response: "update success"})
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
