package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alireza-akbarzadeh/ginflow/pkg/api/handlers"
	"github.com/stretchr/testify/assert"
)

// TestAuthenticationFlow tests the complete authentication flow
func TestAuthenticationFlow(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	t.Run("user registration and login", func(t *testing.T) {
		// Test user registration
		registerReq := handlers.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "Test User",
		}

		w := ts.createRequest("POST", "/api/v1/auth/register", registerReq)
		assert.Equal(t, http.StatusCreated, w.Code)

		var user map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &user)
		assert.NoError(t, err)
		assert.Equal(t, "test@example.com", user["email"])
		assert.Equal(t, "Test User", user["name"])
		assert.NotContains(t, user, "password") // Password should not be exposed
	})

	t.Run("user login", func(t *testing.T) {
		// Test user login
		loginReq := handlers.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		w := ts.createRequest("POST", "/api/v1/auth/login", loginReq)
		assert.Equal(t, http.StatusOK, w.Code)

		var loginResp handlers.LoginResponse
		err := json.Unmarshal(w.Body.Bytes(), &loginResp)
		assert.NoError(t, err)
		assert.NotEmpty(t, loginResp.Token)
		assert.NotNil(t, loginResp.User)
		assert.Equal(t, "test@example.com", loginResp.User.Email)
		assert.Equal(t, "Test User", loginResp.User.Name)
	})

	t.Run("duplicate user registration", func(t *testing.T) {
		// Try to register the same user again
		registerReq := handlers.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			Name:     "Test User",
		}

		w := ts.createRequest("POST", "/api/v1/auth/register", registerReq)
		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("invalid login credentials", func(t *testing.T) {
		// Test invalid password
		loginReq := handlers.LoginRequest{
			Email:    "test@example.com",
			Password: "wrongpassword",
		}

		w := ts.createRequest("POST", "/api/v1/auth/login", loginReq)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("login with non-existent user", func(t *testing.T) {
		// Test non-existent user
		loginReq := handlers.LoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		}

		w := ts.createRequest("POST", "/api/v1/auth/login", loginReq)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("user logout", func(t *testing.T) {
		w := ts.createRequest("POST", "/api/v1/auth/logout", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response, "message")
	})
}

// TestAuthenticationValidation tests input validation for authentication endpoints
func TestAuthenticationValidation(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	t.Run("invalid email format", func(t *testing.T) {
		registerReq := handlers.RegisterRequest{
			Email:    "invalid-email",
			Password: "password123",
			Name:     "Test User",
		}

		w := ts.createRequest("POST", "/api/v1/auth/register", registerReq)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("password too short", func(t *testing.T) {
		registerReq := handlers.RegisterRequest{
			Email:    "test2@example.com",
			Password: "123",
			Name:     "Test User",
		}

		w := ts.createRequest("POST", "/api/v1/auth/register", registerReq)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("name too short", func(t *testing.T) {
		registerReq := handlers.RegisterRequest{
			Email:    "test3@example.com",
			Password: "password123",
			Name:     "A",
		}

		w := ts.createRequest("POST", "/api/v1/auth/register", registerReq)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required fields", func(t *testing.T) {
		w := ts.createRequest("POST", "/api/v1/auth/register", map[string]string{})
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
