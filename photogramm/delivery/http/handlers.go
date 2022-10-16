package http

import (
	"example-restful-api-server/auth"
	"example-restful-api-server/e"
	"example-restful-api-server/models"
	"example-restful-api-server/photogramm"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	useCase photogramm.UseCase
}

type getInput struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type removeInput struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type createInput struct {
	Name string `json:"name"`
}

type responce struct {
	Resp string `json:"responce"`
}

func NewHandler(uc photogramm.UseCase) *Handler {
	return &Handler{
		useCase: uc,
	}
}

func (h *Handler) Upload(c *gin.Context) {
	errStr := "uploadHandler:"

	user := c.MustGet(auth.CtxUserKey).(*models.User)

	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		e.ErrorHandler(c.Writer, http.StatusBadRequest, err, errStr)
		return
	}

	file, fileHeader, err := c.Request.FormFile("photo")
	if err != nil {
		e.ErrorHandler(c.Writer, http.StatusBadRequest, err, errStr)
		return
	}
	defer file.Close()

	if err := isImage(fileHeader); err != nil {
		e.ErrorHandler(c.Writer, http.StatusBadRequest, err, errStr)
		return
	}
	log.Printf("Uploaded File: %+v\n", fileHeader.Filename)
	log.Printf("File Size: %+v\n", fileHeader.Size)
	log.Printf("MIME Header: %+v\n", fileHeader.Header)

	albName := c.Request.FormValue("album_name")

	if err = h.useCase.UploadPhoto(c.Request.Context(), user, albName, file); err != nil {
		e.ErrorHandler(c.Writer, http.StatusBadRequest, err, errStr)
		return
	}
}

func (h *Handler) GetPhoto(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	input := getInput{}

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("getPhoto: %s", err)
		return
	}

	photo, err := h.useCase.GetPhoto(c, user, input.Id)

	if err == e.ErrNotFound {
		log.Printf("getPhoto: %s", err)
		resp := responce{
			Resp: "not found",
		}
		c.JSON(http.StatusOK, resp)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("getPhoto: %s", err)
		return
	}
	c.JSON(http.StatusOK, photo)
}

func (h *Handler) RemovePhoto(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	input := removeInput{}

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("removePhoto: %s", err)
		return
	}

	err := h.useCase.RemovePhoto(c, user, input.Id)
	if err == e.ErrNotFound {
		log.Printf("removePhoto: %s", err)
		resp := responce{Resp: "not found"}
		c.JSON(http.StatusNotFound, resp)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("removePhoto: %s", err)
		return
	}
	resp := responce{Resp: "delete success"}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) CreateAlbum(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	input := createInput{}

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("createAlbum: %s", err)
		return
	}

	err := h.useCase.CreateAlbum(c.Request.Context(), user, input.Name)
	if err == e.ErrAlreadyExist {
		log.Printf("createAlbum: %s", err)
		resp := responce{Resp: "name already in use"}
		c.JSON(http.StatusConflict, resp)
		return
	}

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("createAlbum: %s", err)
		return
	}
	resp := responce{Resp: "success"}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetAlbum(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	input := getInput{}

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("getAlbum: %s", err)
		return
	}

	album, err := h.useCase.GetAlbum(c, user, input.Name)

	if err == e.ErrNotFound {
		log.Printf("getAlbum: %s", err)
		resp := responce{Resp: "not found"}
		c.JSON(http.StatusNotFound, resp)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("getPhoto: %s", err)
		return
	}
	c.JSON(http.StatusOK, album)
}

func (h *Handler) RemoveAlbum(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	input := removeInput{}

	if err := c.BindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("removeAlbum: %s", err)
		return
	}

	err := h.useCase.RemoveAlbum(c, user, input.Name)
	if err == e.ErrNotFound {
		log.Printf("removeAlbum: %s", err)
		resp := responce{Resp: "not found"}
		c.JSON(http.StatusNotFound, resp)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("removePhoto: %s", err)
		return
	}
	resp := responce{Resp: "success"}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetInfo(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	var err error
	user, err = h.useCase.GetInfo(c, user)

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("getInfo: %s", err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func isImage(h *multipart.FileHeader) error {
	s := h.Header.Get("Content-Type")
	if strings.Contains(s, "image") {
		return nil
	}
	return fmt.Errorf("not an image")
}
