package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupProtectedUserRoutes configures protected user routes
func SetupProtectedUserRoutes(router *gin.RouterGroup, handler *handlers.Handler) {
	router.GET("/users", handler.GetAllUsers)
	router.PUT("/users/:id", handler.UpdateUser)
	router.DELETE("/users/:id", handler.DeleteUser)
}
