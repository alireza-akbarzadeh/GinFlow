package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupProtectedBasketRoutes registers protected basket routes
func SetupProtectedBasketRoutes(rg *gin.RouterGroup, handler *handlers.Handler) {
	basket := rg.Group("/basket")
	{
		basket.GET("", handler.GetBasket)
		basket.DELETE("", handler.ClearBasket)
		basket.POST("/items", handler.AddItemToBasket)
		basket.DELETE("/items/:id", handler.RemoveItemFromBasket)
	}
}
