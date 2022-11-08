package http

import (
	"example-restful-api-server/auth"
	e "example-restful-api-server/err"
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

type removeInput struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type createInput struct {
	Name string `json:"name"`
}

type response struct {
	Resp string `json:"response"`
}

type uploadResp struct {
	ID string `json:"id"`
}

func NewHandler(uc photogramm.UseCase) *Handler {
	return &Handler{
		useCase: uc,
	}
}

// Upload
// @Summary      Upload photo
// @Description  Upload photo
// @Tags         api
// @Accept       multipart/form-data
// @Produce      json
// @Security	JWT
// @Param 		photo formData file true "uploaded file data"
// @Param 		album_name formData string false "album name"
// @Success      200  {object}  uploadResp
// @Failure      400
// @Failure      406 {object} response
// @Failure      401
// @Router       /api/photogramm/upload [post]
func (h *Handler) Upload(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}

	file, fileHeader, err := c.Request.FormFile("photo")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer file.Close()

	if err := isImage(fileHeader); err != nil {
		c.JSON(http.StatusNotAcceptable, response{Resp: "not image"})
		return
	}
	log.Printf("Uploaded File: %+v\n", fileHeader.Filename)
	log.Printf("File Size: %+v\n", fileHeader.Size)
	log.Printf("MIME Header: %+v\n", fileHeader.Header)

	albName := c.Request.FormValue("album_name")

	photo := models.NewPhoto(user.Username, user.ID, albName)

	id, err := h.useCase.UploadPhoto(c.Request.Context(), photo, file)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, uploadResp{ID: id})

}

// Get photo
// @Summary      Get photo
// @Description  Get photo by id
// @Tags         api
// @Produce      json
// @Security	JWT
// @Param 		id path string true "photo id"
// @Success      200  {object}  models.Photo
// @Failure      400
// @Failure      404 {object} response
// @Failure      406 {object} response
// @Failure      401
// @Router       /api/photogramm/getphoto/{id} [get]
func (h *Handler) GetPhoto(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	id := c.Param("id")

	photo, err := h.useCase.GetPhoto(c, user, id)

	if err == e.ErrNotFound {
		log.Printf("getPhoto: %s", err)
		resp := response{
			Resp: "not found",
		}
		c.JSON(http.StatusNotFound, resp)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("getPhoto: %s", err)
		return
	}
	c.JSON(http.StatusOK, photo)
}

// Remove photo
// @Summary      Remove photo
// @Description  Remove photo by id
// @Tags         api
// @Accept       json
// @Produce      json
// @Security	JWT
// @Param 		id body removeInput true "photo id"
// @Success      200  {object}  response
// @Failure      400
// @Failure      404 {object} response
// @Failure      406 {object} response
// @Failure      401
// @Router       /api/photogramm/removephoto [delete]
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
		resp := response{Resp: "not found"}
		c.JSON(http.StatusNotFound, resp)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("removePhoto: %s", err)
		return
	}
	resp := response{Resp: "delete success"}
	c.JSON(http.StatusOK, resp)
}

// Create album
// @Summary      Create album
// @Description  Create new album
// @Tags         api
// @Accept       json
// @Produce      json
// @Security	JWT
// @Param 		album_name body createInput true "album name"
// @Success      200  {object}  response
// @Failure      400 {object} response
// @Failure      409 {object} response
// @Failure      401
// @Router       /api/photogramm/createalbum [post]
func (h *Handler) CreateAlbum(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	input := createInput{}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, response{Resp: "incorrect request body"})
		log.Printf("createAlbum: %s", err)
		return
	}

	err := h.useCase.CreateAlbum(c.Request.Context(), user, input.Name)
	if err == e.ErrAlreadyExist {
		log.Printf("createAlbum: %s", err)
		resp := response{Resp: "name already in use"}
		c.JSON(http.StatusConflict, resp)
		return
	}

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("createAlbum: %s", err)
		return
	}
	resp := response{Resp: "success"}
	c.JSON(http.StatusOK, resp)
}

// Get album
// @Summary      Get album
// @Description  Get album by name
// @Tags         api
// @Produce      json
// @Security	JWT
// @Param 		name path string true "album name"
// @Success      200  {object}  models.PhotoAlbum
// @Failure      400
// @Failure      404 {object} response
// @Failure      401
// @Router       /api/photogramm/getalbum/{name} [get]
func (h *Handler) GetAlbum(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	name := c.Param("name")

	album, err := h.useCase.GetAlbum(c, user, name)

	if err == e.ErrNotFound {
		log.Printf("getAlbum: %s", err)
		resp := response{Resp: "not found"}
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

// Remove album
// @Summary      Remove album
// @Description  Remove album by name
// @Tags         api
// @Accept       json
// @Produce      json
// @Security	JWT
// @Param 		name body removeInput true "album name"
// @Success      200  {object}  response
// @Failure      400 {object} response
// @Failure      401
// @Failure      404 {object} response
// @Router       /api/photogramm/removealbum [delete]
func (h *Handler) RemoveAlbum(c *gin.Context) {
	user := c.MustGet(auth.CtxUserKey).(*models.User)

	input := removeInput{}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, &response{Resp: "incorrect request body"})
		log.Printf("removeAlbum: %s", err)
		return
	}

	err := h.useCase.RemoveAlbum(c, user, input.Name)
	if err == e.ErrNotFound {
		log.Printf("removeAlbum: %s", err)
		resp := response{Resp: "not found"}
		c.JSON(http.StatusNotFound, resp)
		return
	}
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		log.Printf("removePhoto: %s", err)
		return
	}
	resp := response{Resp: "success"}
	c.JSON(http.StatusOK, resp)
}

// Get info
// @Summary      Get user info
// @Description  Get user storage summary info
// @Tags         api
// @Produce      json
// @Security	JWT
// @Success      200  {object}  models.User
// @Failure      400
// @Failure      401
// @Router       /api/photogramm/getinfo [get]
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
