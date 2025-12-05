package repository

import (
	"github.com/alireza-akbarzadeh/ginflow/pkg/models"
	"gorm.io/gorm"
)

// CategoryRepository handles category database operations
type CategoryRepository struct {
	DB *gorm.DB
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{DB: db}
}

// Insert creates a new category
func (r *CategoryRepository) Insert(category *models.Category) (*models.Category, error) {
	result := r.DB.Create(category)
	if result.Error != nil {
		return nil, result.Error
	}
	return category, nil
}

// GetAll retrieves all categories
func (r *CategoryRepository) GetAll() ([]*models.Category, error) {
	var categories []*models.Category
	result := r.DB.Find(&categories)
	if result.Error != nil {
		return nil, result.Error
	}
	return categories, nil
}

// Get retrieves a category by ID
func (r *CategoryRepository) Get(id int) (*models.Category, error) {
	var category models.Category
	result := r.DB.First(&category, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &category, nil
}
