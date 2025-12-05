package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/pkg/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupCategoryRoutes configures public category routes
func SetupCategoryRoutes(router *gin.RouterGroup, handler *handlers.Handler) {
	categories := router.Group("/categories")
	{
		categories.GET("", handler.GetAllCategories)
		categories.GET("/:slug", handler.GetCategoryBySlug)
	}
}

// SetupProtectedCategoryRoutes configures protected category routes
func SetupProtectedCategoryRoutes(router *gin.RouterGroup, handler *handlers.Handler) {
	router.POST("/categories", handler.CreateCategory)
}
