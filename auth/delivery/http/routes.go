package http

import (
	"example-restful-api-server/auth"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, uc auth.UseCase) {
	h := NewHandler(uc)

	authRoutes := r.Group("/auth")
	{
		authRoutes.POST("/sign-up", h.SignUp)
		authRoutes.POST("/sign-in", h.SignIn)
	}
}

func RegisterMidRoutes(r *gin.RouterGroup, uc auth.UseCase) {
	h := NewHandler(uc)

	userRoutes := r.Group("/delete")
	{
		userRoutes.DELETE("", h.Delete)
	}

}
