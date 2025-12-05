package repository

import "gorm.io/gorm"

// Models holds all repository models
type Models struct {
	Users      UserRepositoryInterface
	Events     EventRepositoryInterface
	Attendees  AttendeeRepositoryInterface
	Categories CategoryRepositoryInterface
	Comments   CommentRepositoryInterface
	Profiles   ProfileRepositoryInterface
	Products   ProductRepositoryInterface
	Baskets    BasketRepositoryInterface
}

// NewModels creates a new Models instance with all repositories
func NewModels(db *gorm.DB) *Models {
	return &Models{
		Users:      NewUserRepository(db),
		Events:     NewEventRepository(db),
		Attendees:  NewAttendeeRepository(db),
		Categories: NewCategoryRepository(db),
		Comments:   NewCommentRepository(db),
		Profiles:   NewProfileRepository(db),
		Products:   NewProductRepository(db),
		Baskets:    NewBasketRepository(db),
	}
}
