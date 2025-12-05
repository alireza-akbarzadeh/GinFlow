package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupProductRoutes configures product-related routes
func SetupProductRoutes(router *gin.RouterGroup, handler *handlers.Handler) {
	products := router.Group("/products")
	{
		products.GET("", handler.GetAllProducts)
		products.GET("/:id", handler.GetProduct)
		products.GET("/slug/:slug", handler.GetProductBySlug)
		products.GET("/category/:id", handler.GetProductsByCategory)
	}
}

// SetupProtectedProductRoutes configures protected product-related routes
func SetupProtectedProductRoutes(router *gin.RouterGroup, handler *handlers.Handler) {
	products := router.Group("/products")
	{
		products.POST("", handler.CreateProduct)
		products.PUT("/:id", handler.UpdateProduct)
		products.DELETE("/:id", handler.DeleteProduct)
	}
}
