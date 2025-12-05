package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupAttendeeRoutes configures public attendee routes
func SetupAttendeeRoutes(router *gin.RouterGroup, handler *handlers.Handler) {
	attendees := router.Group("/attendees")
	{
		attendees.GET("/:id/events", handler.GetEventsByAttendee)
	}
}
