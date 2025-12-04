package main

import (
	"net/http"
	"strconv"

	"github.com/alireza-akbarzadeh/restful-app/internal/database"
	"github.com/alireza-akbarzadeh/restful-app/internal/messages"
	"github.com/gin-gonic/gin"
)

// @Summary      Create a new event
// @Description  Create a new event (requires authentication)
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        event  body      database.Event  true  "Event object"
// @Success      201    {object}  database.Event
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/events [post]
func (app *application) createEvent(c *gin.Context) {
	var event database.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := app.getUserFromContext(c)
	event.OwnerId = user.Id
	createdEvent, err := app.models.Events.Insert(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrCreateEventFailed})
		return
	}
	c.JSON(http.StatusCreated, createdEvent)
}

// @Summary      Get a single event
// @Description  Get event by ID
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Event ID"
// @Success      200  {object}  database.Event
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/events/{id} [get]
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

// @Summary      Get all events
// @Description  Get a list of all events
// @Tags         Events
// @Accept       json
// @Produce      json
// @Success      200  {array}   database.Event
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/events [get]
func (app *application) getAllEvent(c *gin.Context) {

	event, err := app.models.Events.GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrRetrieveEventFailed})
		return
	}

	c.JSON(http.StatusOK, event)
}

// @Summary      Update an event
// @Description  Update an existing event (requires authentication and ownership)
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id     path      int             true  "Event ID"
// @Param        event  body      database.Event  true  "Event object"
// @Success      200    {object}  database.Event
// @Failure      400    {object}  map[string]string
// @Failure      403    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/events/{id} [put]
func (app *application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidEventID})
		return
	}
	user := app.getUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)
	if existingEvent == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrRetrieveEventFailed})
		return
	}
	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to update this event"})
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

// @Summary      Delete an event
// @Description  Delete an existing event (requires authentication and ownership)
// @Tags         Events
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Event ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/events/{id} [delete]
func (app *application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": messages.ErrInvalidEventID})
		return
	}
	user := app.getUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "faield to retrieve event"})
		return
	}
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrEventNotFound})
		return
	}
	if err := app.models.Events.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrDeleteEventFailed})
		return
	}
	if existingEvent.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this event"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Summary      Add attendee to event
// @Description  Add an attendee to an event (requires authentication and ownership)
// @Tags         Attendees
// @Accept       json
// @Produce      json
// @Param        id      path      int  true  "Event ID"
// @Param        userId  path      int  true  "User ID"
// @Success      201     {object}  database.Attendee
// @Failure      400     {object}  map[string]string
// @Failure      403     {object}  map[string]string
// @Failure      404     {object}  map[string]string
// @Failure      409     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/events/{id}/attendees/{userId} [post]
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
	user := app.getUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to add attendees to this event"})
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

// @Summary      Get attendees for event
// @Description  Get all attendees for a specific event
// @Tags         Attendees
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Event ID"
// @Success      200 {array}   database.User
// @Failure      400 {object}  map[string]string
// @Failure      500 {object}  map[string]string
// @Router       /api/v1/events/{id}/attendees [get]
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

// @Summary      Remove attendee from event
// @Description  Remove an attendee from an event (requires authentication and ownership)
// @Tags         Attendees
// @Accept       json
// @Produce      json
// @Param        id      path  int  true  "Event ID"
// @Param        userId  path  int  true  "User ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      403  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     BearerAuth
// @Router       /api/v1/events/{id}/attendees/{userId} [delete]
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
	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": messages.ErrRetrieveEventFailed})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": messages.ErrEventNotFound})
		return
	}
	user := app.getUserFromContext(c)
	if event.OwnerId != user.Id {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "You are not authorized to delete this event"})
		return
	}

	err = app.models.Attendees.Delete(userId, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete attendee for event"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// @Summary      Get events by attendee
// @Description  Get all events for a specific attendee
// @Tags         Attendees
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Attendee ID"
// @Success      200 {array}   database.Event
// @Failure      400 {object}  map[string]string
// @Failure      500 {object}  map[string]string
// @Router       /api/v1/attendees/{id}/events [get]
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
