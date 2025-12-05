package repository

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
)

type UserRepositoryInterface interface {
	Insert(ctx context.Context, user *models.User) (*models.User, error)
	Get(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetById(ctx context.Context, id int) (*models.User, error)
	UpdatePassword(ctx context.Context, userID int, hashedPassword string) error
	GetAll(ctx context.Context) ([]*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int) error
	UpdateLastLogin(ctx context.Context, id int) error
}

type EventRepositoryInterface interface {
	Insert(ctx context.Context, event *models.Event) (*models.Event, error)
	Get(ctx context.Context, id int) (*models.Event, error)
	GetAll(ctx context.Context) ([]*models.Event, error)
	Update(ctx context.Context, event *models.Event) error
	Delete(ctx context.Context, id int) error
}

type AttendeeRepositoryInterface interface {
	Insert(ctx context.Context, attendee *models.Attendee) (*models.Attendee, error)
	GetByEventAndUser(ctx context.Context, eventID, userID int) (*models.Attendee, error)
	GetByEventAndAttendee(ctx context.Context, eventID, userID int) (*models.Attendee, error)
	GetAttendeesByEvent(ctx context.Context, eventID int) ([]*models.User, error)
	GetEventsByAttendee(ctx context.Context, userID int) ([]*models.Event, error)
	GetEventByAttendee(ctx context.Context, userID int) ([]*models.Event, error)
	Delete(ctx context.Context, userID, eventID int) error
	DeleteByEvent(ctx context.Context, eventID int) error
	DeleteByUser(ctx context.Context, userID int) error
}

type CategoryRepositoryInterface interface {
	Insert(ctx context.Context, category *models.Category) (*models.Category, error)
	GetAll(ctx context.Context) ([]*models.Category, error)
	Get(ctx context.Context, id int) (*models.Category, error)
	GetBySlug(ctx context.Context, slug string) (*models.Category, error)
}

type CommentRepositoryInterface interface {
	Insert(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetByEvent(ctx context.Context, eventID int) ([]*models.Comment, error)
	Delete(ctx context.Context, id int) error
	Get(ctx context.Context, id int) (*models.Comment, error)
}

type ProfileRepositoryInterface interface {
	Insert(ctx context.Context, profile *models.Profile) (*models.Profile, error)
	GetByUserID(ctx context.Context, userID int) (*models.Profile, error)
	GetByUserIDWithUser(ctx context.Context, userID int) (*models.Profile, error)
	Update(ctx context.Context, profile *models.Profile) error
	UpdateByUserID(ctx context.Context, userID int, updates map[string]interface{}) error
	DeleteByUserID(ctx context.Context, id int) error
}

type ProductRepositoryInterface interface {
	Insert(ctx context.Context, product *models.Product) (*models.Product, error)
	GetAll(ctx context.Context, page, limit int, search string, categoryID int) ([]models.Product, int64, error)
	Get(ctx context.Context, id int) (*models.Product, error)
	GetBySlug(ctx context.Context, slug string) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
	GetByUser(ctx context.Context, userID int) ([]models.Product, error)
	GetByCategory(ctx context.Context, categoryID int) ([]models.Product, error)
}

type BasketRepositoryInterface interface {
	GetActiveBasket(ctx context.Context, userID int) (*models.Basket, error)
	CreateBasket(ctx context.Context, basket *models.Basket) error
	AddItem(ctx context.Context, basketID int, item *models.BasketItem) error
	UpdateItemQuantity(ctx context.Context, itemID int, quantity int) error
	RemoveItem(ctx context.Context, itemID int) error
	ClearBasket(ctx context.Context, basketID int) error
}
