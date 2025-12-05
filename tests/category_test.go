package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/alireza-akbarzadeh/ginflow/pkg/models"
	"github.com/stretchr/testify/assert"
)

// TestCategoryManagement tests category CRUD operations
func TestCategoryManagement(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	// Create test user for authentication
	token, _ := ts.createTestUser(t, "categoryuser@example.com", "password123", "Category User")

	t.Run("create category", func(t *testing.T) {
		category := models.Category{
			Name:        "Technology",
			Description: "Technology-related events",
		}

		w := ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, category)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdCategory models.Category
		err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
		assert.NoError(t, err)
		assert.NotZero(t, createdCategory.ID)
		assert.Equal(t, category.Name, createdCategory.Name)
		assert.Equal(t, category.Description, createdCategory.Description)
	})

	t.Run("get all categories", func(t *testing.T) {
		// First create a few categories
		cat1 := models.Category{Name: "Sports", Description: "Sports events"}
		ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, cat1)

		cat2 := models.Category{Name: "Music", Description: "Music concerts and events"}
		ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, cat2)

		// Get all categories
		w := ts.createRequest("GET", "/api/v1/categories", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var categories []models.Category
		err := json.Unmarshal(w.Body.Bytes(), &categories)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(categories), 2)

		// Check that categories have correct data
		foundSports := false
		foundMusic := false
		for _, c := range categories {
			if c.Name == "Sports" {
				assert.Equal(t, "Sports events", c.Description)
				foundSports = true
			}
			if c.Name == "Music" {
				assert.Equal(t, "Music concerts and events", c.Description)
				foundMusic = true
			}
		}
		assert.True(t, foundSports, "Sports category not found")
		assert.True(t, foundMusic, "Music category not found")
	})

	t.Run("category validation", func(t *testing.T) {
		// Test empty name
		invalidCategory := models.Category{Name: "", Description: "Valid description"}
		w := ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, invalidCategory)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Test name too short
		shortNameCategory := models.Category{Name: "AB", Description: "Valid description"}
		w = ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, shortNameCategory)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Test duplicate name
		duplicateCategory := models.Category{Name: "Technology", Description: "Duplicate technology category"}
		w = ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, duplicateCategory)
		assert.Equal(t, http.StatusInternalServerError, w.Code) // GORM returns 500 for unique constraint violations
	})

	t.Run("get categories without authentication", func(t *testing.T) {
		// Getting categories should work without authentication
		w := ts.createRequest("GET", "/api/v1/categories", nil)
		assert.Equal(t, http.StatusOK, w.Code)

		var categories []models.Category
		err := json.Unmarshal(w.Body.Bytes(), &categories)
		assert.NoError(t, err)
		assert.IsType(t, []models.Category{}, categories)
	})

	t.Run("create category requires authentication", func(t *testing.T) {
		category := models.Category{
			Name:        "Unauthenticated Category",
			Description: "This should fail",
		}

		w := ts.createRequest("POST", "/api/v1/categories", category)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

// TestCategoryEdgeCases tests edge cases for categories
func TestCategoryEdgeCases(t *testing.T) {
	ts := SetupTestSuite(t)
	if ts == nil {
		t.Skip("Test suite setup failed")
		return
	}
	defer ts.TeardownTestSuite(t)

	token, _ := ts.createTestUser(t, "edgecategoryuser@example.com", "password123", "Edge Category User")

	t.Run("category with empty description", func(t *testing.T) {
		category := models.Category{
			Name:        "Empty Description Category",
			Description: "",
		}

		w := ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, category)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdCategory models.Category
		err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
		assert.NoError(t, err)
		assert.Equal(t, category.Name, createdCategory.Name)
		assert.Equal(t, "", createdCategory.Description)
	})

	t.Run("category name with special characters", func(t *testing.T) {
		category := models.Category{
			Name:        "Tech & Innovation",
			Description: "Category with special characters",
		}

		w := ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, category)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdCategory models.Category
		err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
		assert.NoError(t, err)
		assert.Equal(t, "Tech & Innovation", createdCategory.Name)
	})

	t.Run("very long category name", func(t *testing.T) {
		longName := "This is a very long category name that exceeds normal length but should still be accepted since we don't have explicit length limits in validation"
		category := models.Category{
			Name:        longName,
			Description: "Category with very long name",
		}

		w := ts.createAuthenticatedRequest("POST", "/api/v1/categories", token, category)
		assert.Equal(t, http.StatusCreated, w.Code)

		var createdCategory models.Category
		err := json.Unmarshal(w.Body.Bytes(), &createdCategory)
		assert.NoError(t, err)
		assert.Equal(t, longName, createdCategory.Name)
	})
}
