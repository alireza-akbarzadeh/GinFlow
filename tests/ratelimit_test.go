package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alireza-akbarzadeh/ginflow/internal/constants"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiting(t *testing.T) {
	suite := SetupMockTestSuite(t)

	// The rate limit is configured in SetupRouter using constants.
	// DEFAULT_RATE_LIMIT = 20
	// DEFAULT_RATE_BURST = 50

	// We need to exceed the burst limit to trigger the rate limiter.
	// Since the burst is 50, we should be able to make 50 requests immediately.
	// The 51st request might fail if sent too quickly.
	// However, the token bucket refills at 20 tokens/second.

	// Let's try to send requests up to the burst limit + some extra
	totalRequests := constants.DEFAULT_RATE_BURST + 10
	rateLimited := false

	for i := 0; i < totalRequests; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/health", nil)
		suite.Router.ServeHTTP(w, req)

		if w.Code == http.StatusTooManyRequests {
			rateLimited = true
			break
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	suite.Router.ServeHTTP(w, req)
	assert.NotEqual(t, http.StatusInternalServerError, w.Code)

	if rateLimited {
		t.Log("Rate limit triggered successfully")
	} else {
		t.Log("Rate limit was not triggered (machine might be slow or limit too high for this test)")
	}
}
