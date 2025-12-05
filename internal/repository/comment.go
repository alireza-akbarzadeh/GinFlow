package repository

import (
	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"gorm.io/gorm"
)

// CommentRepository handles comment database operations
type CommentRepository struct {
	DB *gorm.DB
}

// NewCommentRepository creates a new CommentRepository
func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

// Insert creates a new comment
func (r *CommentRepository) Insert(comment *models.Comment) (*models.Comment, error) {
	result := r.DB.Create(comment)
	if result.Error != nil {
		return nil, result.Error
	}
	// Preload user info to return complete object
	r.DB.Preload("User").First(comment, comment.ID)
	return comment, nil
}

// GetByEvent retrieves all comments for a specific event
func (r *CommentRepository) GetByEvent(eventID int) ([]*models.Comment, error) {
	var comments []*models.Comment
	result := r.DB.Where("event_id = ?", eventID).Preload("User").Order("created_at desc").Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}
	return comments, nil
}

// Delete removes a comment
func (r *CommentRepository) Delete(id int) error {
	result := r.DB.Delete(&models.Comment{}, id)
	return result.Error
}

// Get retrieves a comment by ID
func (r *CommentRepository) Get(id int) (*models.Comment, error) {
	var comment models.Comment
	result := r.DB.First(&comment, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &comment, nil
}
