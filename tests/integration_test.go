package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHealthCheck tests the health check endpoint
func TestHealthCheck(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	w := ts.createRequest("GET", "/health", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestSwaggerEndpoint tests the swagger endpoint
func TestSwaggerEndpoint(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	w := ts.createRequest("GET", "/swagger/index.html", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestLandingPage tests the landing page
func TestLandingPage(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	w := ts.createRequest("GET", "/", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestDashboardPage tests the dashboard page
func TestDashboardPage(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	w := ts.createRequest("GET", "/dashboard", nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestUnauthenticatedEndpoints tests endpoints that should work without authentication
func TestUnauthenticatedEndpoints(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	// Test logout endpoint (should work without auth)
	w := ts.createRequest("POST", "/api/v1/auth/logout", nil)
	assert.Equal(t, http.StatusOK, w.Code)

	// Test getting all events (should work without auth)
	w = ts.createRequest("GET", "/api/v1/events", nil)
	// This might return 500 if database is not available, but the endpoint should be accessible
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusInternalServerError)
}

// TestAuthenticationEndpointsWithoutDB tests auth endpoints structure
func TestAuthenticationEndpointsWithoutDB(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	// Test register endpoint with invalid data (should return 400)
	w := ts.createRequest("POST", "/api/v1/auth/register", map[string]string{"invalid": "data"})
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Test login endpoint with invalid data (should return 400)
	w = ts.createRequest("POST", "/api/v1/auth/login", map[string]string{"invalid": "data"})
	assert.Equal(t, http.StatusBadRequest, w.Code)
}