package tests

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/alireza-akbarzadeh/ginflow/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestCommentManagement tests comment CRUD operations
func TestCommentManagement(t *testing.T) {
	ts := SetupMockTestSuite(t)

	// Create test users
	user1ID := 1
	user2ID := 2
	token1, _ := ts.GenerateToken(user1ID)
	// token2 is not used in this test suite as we mock the repo responses directly

	mockUserRepo := ts.Mocks.Users.(*mocks.UserRepositoryMock)
	mockUserRepo.On("Get", mock.Anything, user1ID).Return(&models.User{ID: user1ID, Email: "commentuser1@example.com"}, nil)
	mockUserRepo.On("Get", mock.Anything, user2ID).Return(&models.User{ID: user2ID, Email: "commentuser2@example.com"}, nil)

	eventID := 1
	event := &models.Event{
		ID:          eventID,
		Name:        "Event for Comments",
		Description: "This event will have comments",
		Date:        "2025-12-31",
		Location:    "Comment Location",
	}

	mockEventRepo := ts.Mocks.Events.(*mocks.EventRepositoryMock)
	mockCommentRepo := ts.Mocks.Comments.(*mocks.CommentRepositoryMock)

	t.Run("create comment", func(t *testing.T) {
		comment := models.Comment{
			Content: "This is a test comment",
		}

		mockEventRepo.On("Get", mock.Anything, eventID).Return(event, nil).Once()
		mockCommentRepo.On("Insert", mock.Anything, mock.MatchedBy(func(c *models.Comment) bool {
			return c.Content == comment.Content && c.UserID == user1ID && c.EventID == eventID
		})).Return(&models.Comment{ID: 1, Content: comment.Content, UserID: user1ID, EventID: eventID}, nil).Once()

		w := ts.createAuthenticatedRequest("POST", "/api/v1/events/"+strconv.Itoa(eventID)+"/comments", token1, comment)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdComment models.Comment
		err := json.Unmarshal(w.Body.Bytes(), &createdComment)
		assert.NoError(t, err)
		assert.Equal(t, comment.Content, createdComment.Content)
	})

	t.Run("get event comments", func(t *testing.T) {
		comments := []*models.Comment{
			{ID: 1, Content: "First comment", UserID: user1ID, EventID: eventID},
			{ID: 2, Content: "Second comment", UserID: user2ID, EventID: eventID},
		}

		mockEventRepo.On("Get", mock.Anything, eventID).Return(event, nil).Once()
		mockCommentRepo.On("GetByEvent", mock.Anything, eventID).Return(comments, nil).Once()

		// Get all comments for the event
		w := ts.createRequest("GET", "/api/v1/events/"+strconv.Itoa(eventID)+"/comments", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var respComments []models.Comment
		err := json.Unmarshal(w.Body.Bytes(), &respComments)
		assert.NoError(t, err)
		assert.Equal(t, len(comments), len(respComments))
	})

	t.Run("delete own comment", func(t *testing.T) {
		commentID := 1
		comment := &models.Comment{ID: commentID, Content: "Comment to delete", UserID: user1ID, EventID: eventID}

		mockEventRepo.On("Get", mock.Anything, eventID).Return(event, nil).Once()
		mockCommentRepo.On("Get", mock.Anything, commentID).Return(comment, nil).Once()
		mockCommentRepo.On("Delete", mock.Anything, commentID).Return(nil).Once()

		// Delete the comment
		w := ts.createAuthenticatedRequest("DELETE", "/api/v1/events/"+strconv.Itoa(eventID)+"/comments/"+strconv.Itoa(commentID), token1, nil)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("comment validation", func(t *testing.T) {
		// Test empty content
		invalidComment := models.Comment{Content: ""}
		w := ts.createAuthenticatedRequest("POST", "/api/v1/events/"+strconv.Itoa(eventID)+"/comments", token1, invalidComment)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Test comment on non-existent event
		nonExistentEventID := 99999
		validComment := models.Comment{Content: "Valid comment"}

		mockEventRepo.On("Get", mock.Anything, nonExistentEventID).Return(nil, nil).Once()

		w = ts.createAuthenticatedRequest("POST", "/api/v1/events/"+strconv.Itoa(nonExistentEventID)+"/comments", token1, validComment)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestCommentAuthorization tests comment authorization
func TestCommentAuthorization(t *testing.T) {
	ts := SetupMockTestSuite(t)

	user1ID := 1
	user2ID := 2
	// token1 is not used
	token2, _ := ts.GenerateToken(user2ID)

	mockUserRepo := ts.Mocks.Users.(*mocks.UserRepositoryMock)
	mockUserRepo.On("Get", mock.Anything, user1ID).Return(&models.User{ID: user1ID, Email: "authcommentuser1@example.com"}, nil)
	mockUserRepo.On("Get", mock.Anything, user2ID).Return(&models.User{ID: user2ID, Email: "authcommentuser2@example.com"}, nil)

	eventID := 1
	event := &models.Event{
		ID:          eventID,
		Name:        "Event for Auth Test",
		Description: "Testing comment authorization",
		Date:        "2025-12-31",
		Location:    "Auth Location",
	}

	mockEventRepo := ts.Mocks.Events.(*mocks.EventRepositoryMock)
	mockCommentRepo := ts.Mocks.Comments.(*mocks.CommentRepositoryMock)

	t.Run("cannot delete others comment", func(t *testing.T) {
		commentID := 1
		comment := &models.Comment{ID: commentID, Content: "User 1's comment", UserID: user1ID, EventID: eventID}

		mockEventRepo.On("Get", mock.Anything, eventID).Return(event, nil).Once()
		mockCommentRepo.On("Get", mock.Anything, commentID).Return(comment, nil).Once()

		// User 2 tries to delete User 1's comment
		w := ts.createAuthenticatedRequest("DELETE", "/api/v1/events/"+strconv.Itoa(eventID)+"/comments/"+strconv.Itoa(commentID), token2, nil)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("unauthenticated cannot create comment", func(t *testing.T) {
		comment := models.Comment{Content: "Unauthenticated comment"}
		w := ts.createRequest("POST", "/api/v1/events/"+strconv.Itoa(eventID)+"/comments", comment)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
