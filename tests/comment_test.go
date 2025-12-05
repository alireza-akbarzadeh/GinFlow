package tests

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	models "github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestCommentManagement tests comment CRUD operations
func TestCommentManagement(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	// Create test users and event
	token1, user1 := ts.createTestUser(t, "commentuser1@example.com", "password123", "Comment User 1")
	token2, user2 := ts.createTestUser(t, "commentuser2@example.com", "password123", "Comment User 2")

	event := models.Event{
		Name:        "Event for Comments",
		Description: "This event will have comments",
		Date:        "2025-12-31",
		Location:    "Comment Location",
	}
	createdEvent := ts.createTestEvent(t, token1, event)

	t.Run("create comment", func(t *testing.T) {
		comment := models.Comment{
			Content: "This is a test comment",
		}

		w := ts.createAuthenticatedRequest("POST", "/api/v1/events/"+strconv.Itoa(createdEvent.ID)+"/comments", token1, comment)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdComment models.Comment
		err := json.Unmarshal(w.Body.Bytes(), &createdComment)
		assert.NoError(t, err)
		assert.NotZero(t, createdComment.ID)
		assert.Equal(t, comment.Content, createdComment.Content)
		assert.Equal(t, user1.ID, createdComment.UserID)
		assert.Equal(t, createdEvent.ID, createdComment.EventID)
		assert.NotZero(t, createdComment.CreatedAt)
	})

	t.Run("get event comments", func(t *testing.T) {
		// First add a few comments
		comment1 := models.Comment{Content: "First comment"}
		ts.createAuthenticatedRequest("POST", "/api/v1/events/"+strconv.Itoa(createdEvent.ID)+"/comments", token1, comment1)

		comment2 := models.Comment{Content: "Second comment"}
		ts.createAuthenticatedRequest("POST", "/api/v1/events/"+strconv.Itoa(createdEvent.ID)+"/comments", token2, comment2)

		// Get all comments for the event
		w := ts.createRequest("GET", "/api/v1/events/"+strconv.Itoa(createdEvent.ID)+"/comments", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var comments []models.Comment
		err := json.Unmarshal(w.Body.Bytes(), &comments)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(comments), 2)

		// Check that comments have correct data
		foundComment1 := false
		foundComment2 := false
		for _, c := range comments {
			if c.Content == "First comment" {
				assert.Equal(t, user1.ID, c.UserID)
				foundComment1 = true
			}
			if c.Content == "Second comment" {
				assert.Equal(t, user2.ID, c.UserID)
				foundComment2 = true
			}
		}
		assert.True(t, foundComment1, "First comment not found")
		assert.True(t, foundComment2, "Second comment not found")
	})

	t.Run("delete own comment", func(t *testing.T) {
		// Create a comment to delete
		comment := models.Comment{Content: "Comment to delete"}
		createResp := ts.createAuthenticatedRequest("POST", "/api/v1/events/"+strconv.Itoa(createdEvent.ID)+"/comments", token1, comment)
		assert.Equal(t, http.StatusCreated, createResp.Code)

		var createdComment models.Comment
		err := json.Unmarshal(createResp.Body.Bytes(), &createdComment)
		assert.NoError(t, err)

		// Delete the comment
		w := ts.createAuthenticatedRequest("DELETE", "/api/v1/events/"+strconv.Itoa(createdEvent.ID)+"/comments/"+strconv.Itoa(createdComment.ID), token1, nil)
		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify comment is deleted
		getResp := ts.createRequest("GET", "/api/v1/events/"+strconv.Itoa(createdEvent.ID)+"/comments", nil)
		assert.Equal(t, http.StatusOK, getResp.Code)

		var comments []models.Comment
		err = json.Unmarshal(getResp.Body.Bytes(), &comments)
		assert.NoError(t, err)

		// Check that the deleted comment is not in the list
		for _, c := range comments {
			assert.NotEqual(t, createdComment.ID, c.ID, "Deleted comment still exists")
		}
	})

	t.Run("comment validation", func(t *testing.T) {
		// Test empty content
		invalidComment := models.Comment{Content: ""}
		w := ts.createAuthenticatedRequest("POST", "/api/v1/events/"+strconv.Itoa(createdEvent.ID)+"/comments", token1, invalidComment)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Test comment on non-existent event
		validComment := models.Comment{Content: "Valid comment"}
		w = ts.createAuthenticatedRequest("POST", "/api/v1/events/99999/comments", token1, validComment)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestCommentAuthorization tests comment authorization
func TestCommentAuthorization(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	// Create test users and event
	token1, _ := ts.createTestUser(t, "authcommentuser1@example.com", "password123", "Auth Comment User 1")
	token2, _ := ts.createTestUser(t, "authcommentuser2@example.com", "password123", "Auth Comment User 2")

	event := models.Event{
		Name:        "Event for Auth Test",
		Description: "Testing comment authorization",
		Date:        "2025-12-31",
		Location:    "Auth Location",
	}
	createdEvent := ts.createTestEvent(t, token1, event)

	t.Run("cannot delete others comment", func(t *testing.T) {
		// User 1 creates a comment
		comment := models.Comment{Content: "User 1's comment"}
		createResp := ts.createAuthenticatedRequest("POST", "/api/v1/events/"+strconv.Itoa(createdEvent.ID)+"/comments", token1, comment)
		assert.Equal(t, http.StatusCreated, createResp.Code)

		var createdComment models.Comment
		err := json.Unmarshal(createResp.Body.Bytes(), &createdComment)
		assert.NoError(t, err)

		// User 2 tries to delete User 1's comment
		w := ts.createAuthenticatedRequest("DELETE", "/api/v1/events/"+strconv.Itoa(createdEvent.ID)+"/comments/"+strconv.Itoa(createdComment.ID), token2, nil)
		assert.Equal(t, http.StatusForbidden, w.Code)

		// Verify comment still exists
		getResp := ts.createRequest("GET", "/api/v1/events/"+strconv.Itoa(createdEvent.ID)+"/comments", nil)
		assert.Equal(t, http.StatusOK, getResp.Code)

		var comments []models.Comment
		err = json.Unmarshal(getResp.Body.Bytes(), &comments)
		assert.NoError(t, err)

		found := false
		for _, c := range comments {
			if c.ID == createdComment.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Comment was incorrectly deleted")
	})

	t.Run("unauthenticated cannot create comment", func(t *testing.T) {
		comment := models.Comment{Content: "Unauthenticated comment"}
		w := ts.createRequest("POST", "/api/v1/events/"+strconv.Itoa(createdEvent.ID)+"/comments", comment)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
