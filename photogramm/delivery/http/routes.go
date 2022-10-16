package http

import (
	"example-restful-api-server/photogramm"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, uc photogramm.UseCase) {
	h := NewHandler(uc)

	photoRoutes := r.Group("/photogramm")
	{
		photoRoutes.POST("upload", h.Upload)
		photoRoutes.GET("getphoto", h.GetPhoto)
		photoRoutes.DELETE("removephoto", h.RemovePhoto)
		photoRoutes.POST("createalbum", h.CreateAlbum)
		photoRoutes.GET("getalbum", h.GetAlbum)
		photoRoutes.DELETE("removealbum", h.RemoveAlbum)
		photoRoutes.GET("getinfo", h.GetInfo)

	}
}
