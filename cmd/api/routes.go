package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	g := gin.Default()

	v1 := g.Group("/api/v1")
	{
		v1.GET("/events", app.getAllEvent)
		v1.GET("/events/:id", app.getEvent)
		v1.GET("/events/:id/attendees", app.getAttendeeForEvent)
		v1.GET("/attendees/:id/events", app.getEventByAttendee)
		v1.POST("/auth/register", app.registerUser)
		v1.POST("/auth/login", app.login)
	}
	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	{
		authGroup.PUT("/events/:id", app.updateEvent)
		authGroup.DELETE("/events/:id", app.deleteEvent)
		authGroup.POST("/events", app.createEvent)
		authGroup.POST("/events/:id/attendees/:userId", app.addAttendeeToEvent)
		authGroup.DELETE("/events/:id/attendees/:userId", app.deleteAttendeeForEvent)

	}
	return g
}
