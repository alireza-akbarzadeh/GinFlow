package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes configures authentication-related routes
func SetupAuthRoutes(router *gin.RouterGroup, handler *handlers.Handler) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", handler.Register)
		auth.POST("/login", handler.Login)
		auth.POST("/logout", handler.Logout)
	}
}

// SetupProtectedAuthRoutes configures protected authentication-related routes
func SetupProtectedAuthRoutes(router *gin.RouterGroup, handler *handlers.Handler) {
	auth := router.Group("/auth")
	{
		auth.PUT("/password", handler.UpdatePassword)
	}
}
