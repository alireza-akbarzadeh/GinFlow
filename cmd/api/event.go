package main

import (
	"net/http"
	"strconv"

	"github.com/alireza-akbarzadeh/restful-app/internal/database"
	"github.com/alireza-akbarzadeh/restful-app/internal/messages"
	"github.com/gin-gonic/gin"
)

func (app *application) createEvent(c *gin.Context) {
	var event database.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdEvent, err := app.models.Events.Insert(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrCreateEventFailed})
		return
	}
	c.JSON(http.StatusCreated, createdEvent)
}

func (app *application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidEventID})
		return
	}

	event, err := app.models.Events.Get(id)
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrEventNotFound})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrRetrieveEventFailed})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (app *application) getAllEvent(c *gin.Context) {

	event, err := app.models.Events.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrRetrieveEventFailed})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidEventID})
		return
	}
	existingEvent, err := app.models.Events.Get(id)
	if existingEvent == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrRetrieveEventFailed})
		return
	}

	updatedEvent := &database.Event{}
	if err := c.ShouldBindJSON(updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedEvent.Id = id
	if err := app.models.Events.Update(updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrUpdateEventFailed})
		return
	}

	c.JSON(http.StatusOK, updatedEvent)
}

func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidEventID})
		return
	}
	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrDeleteEventFailed})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (app *application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidEventID})
		return
	}
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidUserId})
	}
	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrRetrieveEventFailed})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrEventNotFound})
		return
	}
	userToAdd, err := app.models.Users.GetById(userId)
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrUserNotFound})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.Id, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrRetrieveAttendeeFailed})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": messages.ErrAttendeeExist})
		return
	}
	attendee := database.Attendee{
		EventId: event.Id,
		UserId:  userToAdd.Id,
	}
	_, err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}
	c.JSON(http.StatusCreated, attendee)
}

func (app *application) getAttendeeForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}
	users, err := app.models.Attendees.GetAttendeesByEvent(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve attendees for event"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (app *application) deleteAttendeeForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidEventID})
		return
	}
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidUserId})
		return
	}
	err = app.models.Attendees.Delete(userId, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete attendee for event"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func (app *application) getEventByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}
	event, err := app.models.Attendees.GetEventByAttendee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrRetrieveEventFailed})
		return
	}
	c.JSON(http.StatusOK, event)
}
