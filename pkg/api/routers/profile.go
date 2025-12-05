package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/pkg/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupProtectedProfileRoutes configures protected profile routes
func SetupProtectedProfileRoutes(router *gin.RouterGroup, handler *handlers.Handler) {
	router.GET("/profile", handler.GetProfile)
	router.POST("/profile", handler.CreateProfile)
	router.PUT("/profile", handler.UpdateProfile)
	router.DELETE("/profile", handler.DeleteProfile)
}
