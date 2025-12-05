package routers

import (
	"github.com/alireza-akbarzadeh/ginflow/pkg/api/handlers"
	"github.com/gin-gonic/gin"
)

// SetupEventRoutes configures public event routes
func SetupEventRoutes(router *gin.RouterGroup, handler *handlers.Handler) {
	events := router.Group("/events")
	{
		events.GET("", handler.GetAllEvents)
		events.GET("/:id", handler.GetEvent)
		events.GET("/:id/attendees", handler.GetAttendees)
		events.GET("/:id/comments", handler.GetEventComments)
	}
}

// SetupProtectedEventRoutes configures protected event routes
func SetupProtectedEventRoutes(router *gin.RouterGroup, handler *handlers.Handler) {
	// Event management
	router.POST("/events", handler.CreateEvent)
	router.PUT("/events/:id", handler.UpdateEvent)
	router.DELETE("/events/:id", handler.DeleteEvent)

	// Comment management
	router.POST("/events/:id/comments", handler.CreateComment)
	router.DELETE("/events/:id/comments/:commentId", handler.DeleteComment)

	// Attendee management
	router.POST("/events/:id/attendees/:userId", handler.AddAttendee)
	router.DELETE("/events/:id/attendees/:userId", handler.RemoveAttendee)
}
