package tests

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	models "github.com/alireza-akbarzadeh/ginflow/pkg/models"
	"github.com/stretchr/testify/assert"
)

// TestEventManagement tests the complete event management flow
func TestEventManagement(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	// Create a test user and get token
	token, user := ts.createTestUser(t, "eventuser@example.com", "password123", "Event User")

	t.Run("create event", func(t *testing.T) {
		event := models.Event{
			Name:        "Test Event",
			Description: "This is a test event description",
			Date:        "2025-12-31",
			Location:    "Test Location",
		}

		w := ts.createAuthenticatedRequest("POST", "/api/v1/events", token, event)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdEvent models.Event
		err := json.Unmarshal(w.Body.Bytes(), &createdEvent)
		assert.NoError(t, err)
		assert.NotZero(t, createdEvent.ID)
		assert.Equal(t, event.Name, createdEvent.Name)
		assert.Equal(t, event.Description, createdEvent.Description)
		assert.Equal(t, event.Location, createdEvent.Location)
		assert.Equal(t, user.ID, createdEvent.OwnerID)
	})

	t.Run("get all events", func(t *testing.T) {
		w := ts.createRequest("GET", "/api/v1/events", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var events []models.Event
		err := json.Unmarshal(w.Body.Bytes(), &events)
		assert.NoError(t, err)
		assert.NotEmpty(t, events)
		assert.GreaterOrEqual(t, len(events), 1)
	})

	t.Run("get single event", func(t *testing.T) {
		// First create an event to get its ID
		event := models.Event{
			Name:        "Single Event",
			Description: "Description for single event",
			Date:        "2025-12-31",
			Location:    "Single Location",
		}
		createdEvent := ts.createTestEvent(t, token, event)

		// Now get the event by ID
		w := ts.createRequest("GET", "/api/v1/events/"+strconv.Itoa(createdEvent.ID), nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var retrievedEvent models.Event
		err := json.Unmarshal(w.Body.Bytes(), &retrievedEvent)
		assert.NoError(t, err)
		assert.Equal(t, createdEvent.ID, retrievedEvent.ID)
		assert.Equal(t, createdEvent.Name, retrievedEvent.Name)
	})

	t.Run("update event", func(t *testing.T) {
		// Create an event first
		event := models.Event{
			Name:        "Original Event",
			Description: "Original description",
			Date:        "2025-12-31",
			Location:    "Original Location",
		}
		createdEvent := ts.createTestEvent(t, token, event)

		// Update the event
		updatedEvent := models.Event{
			Name:        "Updated Event",
			Description: "Updated description",
			Date:        "2025-12-31",
			Location:    "Updated Location",
		}

		w := ts.createAuthenticatedRequest("PUT", "/api/v1/events/"+strconv.Itoa(createdEvent.ID), token, updatedEvent)
		assert.Equal(t, http.StatusOK, w.Code)

		var result models.Event
		err := json.Unmarshal(w.Body.Bytes(), &result)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Event", result.Name)
		assert.Equal(t, "Updated description", result.Description)
		assert.Equal(t, "Updated Location", result.Location)
	})

	t.Run("delete event", func(t *testing.T) {
		// Create an event first
		event := models.Event{
			Name:        "Event to Delete",
			Description: "This event will be deleted",
			Date:        "2025-12-31",
			Location:    "Delete Location",
		}
		createdEvent := ts.createTestEvent(t, token, event)

		// Delete the event
		w := ts.createAuthenticatedRequest("DELETE", "/api/v1/events/"+strconv.Itoa(createdEvent.ID), token, nil)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Try to get the deleted event (should fail)
		w = ts.createRequest("GET", "/api/v1/events/"+strconv.Itoa(createdEvent.ID), nil)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestEventAuthorization tests that users can only modify their own events
func TestEventAuthorization(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	// Create two users
	token1, _ := ts.createTestUser(t, "user1@example.com", "password123", "User One")
	token2, _ := ts.createTestUser(t, "user2@example.com", "password123", "User Two")

	// User 1 creates an event
	event := models.Event{
		Name:        "User1 Event",
		Description: "Event created by user 1",
		Date:        "2025-12-31",
		Location:    "User1 Location",
	}
	createdEvent := ts.createTestEvent(t, token1, event)

	t.Run("user cannot update another user's event", func(t *testing.T) {
		updatedEvent := models.Event{
			Name:        "Hacked Event",
			Description: "This should not work",
			Date:        "2025-12-31",
			Location:    "Hacked Location",
		}

		w := ts.createAuthenticatedRequest("PUT", "/api/v1/events/"+strconv.Itoa(createdEvent.ID), token2, updatedEvent)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("user cannot delete another user's event", func(t *testing.T) {
		w := ts.createAuthenticatedRequest("DELETE", "/api/v1/events/"+strconv.Itoa(createdEvent.ID), token2, nil)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("owner can update their event", func(t *testing.T) {
		updatedEvent := models.Event{
			Name:        "Updated by Owner",
			Description: "Updated description by owner",
			Date:        "2025-12-31",
			Location:    "Updated Location",
		}

		w := ts.createAuthenticatedRequest("PUT", "/api/v1/events/"+strconv.Itoa(createdEvent.ID), token1, updatedEvent)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestEventValidation tests input validation for events
func TestEventValidation(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	token, _ := ts.createTestUser(t, "validation@example.com", "password123", "Validation User")

	t.Run("missing required fields", func(t *testing.T) {
		w := ts.createAuthenticatedRequest("POST", "/api/v1/events", token, map[string]string{})
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("name too short", func(t *testing.T) {
		event := models.Event{
			Name:        "A",
			Description: "Valid description",
			Date:        "2025-12-31",
			Location:    "Valid location",
		}
		w := ts.createAuthenticatedRequest("POST", "/api/v1/events", token, event)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("description too short", func(t *testing.T) {
		event := models.Event{
			Name:        "Valid Name",
			Description: "Short",
			Date:        "2025-12-31",
			Location:    "Valid location",
		}
		w := ts.createAuthenticatedRequest("POST", "/api/v1/events", token, event)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid date format", func(t *testing.T) {
		event := models.Event{
			Name:        "Valid Name",
			Description: "Valid description that is long enough",
			Date:        "invalid-date",
			Location:    "Valid location",
		}
		w := ts.createAuthenticatedRequest("POST", "/api/v1/events", token, event)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("location too short", func(t *testing.T) {
		event := models.Event{
			Name:        "Valid Name",
			Description: "Valid description that is long enough",
			Date:        "2025-12-31",
			Location:    "A",
		}
		w := ts.createAuthenticatedRequest("POST", "/api/v1/events", token, event)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
