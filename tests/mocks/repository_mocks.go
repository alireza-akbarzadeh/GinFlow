package mocks

import (
	"context"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) Insert(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) Get(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) GetById(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRepositoryMock) UpdatePassword(ctx context.Context, userID int, hashedPassword string) error {
	args := m.Called(ctx, userID, hashedPassword)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetAll(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *UserRepositoryMock) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *UserRepositoryMock) UpdateLastLogin(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type EventRepositoryMock struct {
	mock.Mock
}

func (m *EventRepositoryMock) Insert(ctx context.Context, event *models.Event) (*models.Event, error) {
	args := m.Called(ctx, event)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *EventRepositoryMock) Get(ctx context.Context, id int) (*models.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *EventRepositoryMock) GetAll(ctx context.Context) ([]*models.Event, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *EventRepositoryMock) Update(ctx context.Context, event *models.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *EventRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type AttendeeRepositoryMock struct {
	mock.Mock
}

func (m *AttendeeRepositoryMock) Insert(ctx context.Context, attendee *models.Attendee) (*models.Attendee, error) {
	args := m.Called(ctx, attendee)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attendee), args.Error(1)
}

func (m *AttendeeRepositoryMock) GetByEventAndUser(ctx context.Context, eventID, userID int) (*models.Attendee, error) {
	args := m.Called(ctx, eventID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attendee), args.Error(1)
}

func (m *AttendeeRepositoryMock) GetByEventAndAttendee(ctx context.Context, eventID, userID int) (*models.Attendee, error) {
	args := m.Called(ctx, eventID, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Attendee), args.Error(1)
}

func (m *AttendeeRepositoryMock) GetAttendeesByEvent(ctx context.Context, eventID int) ([]*models.User, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *AttendeeRepositoryMock) GetEventsByAttendee(ctx context.Context, userID int) ([]*models.Event, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *AttendeeRepositoryMock) GetEventByAttendee(ctx context.Context, userID int) ([]*models.Event, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *AttendeeRepositoryMock) Delete(ctx context.Context, userID, eventID int) error {
	args := m.Called(ctx, userID, eventID)
	return args.Error(0)
}

func (m *AttendeeRepositoryMock) DeleteByEvent(ctx context.Context, eventID int) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}

func (m *AttendeeRepositoryMock) DeleteByUser(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type CategoryRepositoryMock struct {
	mock.Mock
}

func (m *CategoryRepositoryMock) Insert(ctx context.Context, category *models.Category) (*models.Category, error) {
	args := m.Called(ctx, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) GetAll(ctx context.Context) ([]*models.Category, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) Get(ctx context.Context, id int) (*models.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Category), args.Error(1)
}

type CommentRepositoryMock struct {
	mock.Mock
}

func (m *CommentRepositoryMock) Insert(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	args := m.Called(ctx, comment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *CommentRepositoryMock) GetByEvent(ctx context.Context, eventID int) ([]*models.Comment, error) {
	args := m.Called(ctx, eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Comment), args.Error(1)
}

func (m *CommentRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *CommentRepositoryMock) Get(ctx context.Context, id int) (*models.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

type ProfileRepositoryMock struct {
	mock.Mock
}

func (m *ProfileRepositoryMock) Insert(ctx context.Context, profile *models.Profile) (*models.Profile, error) {
	args := m.Called(ctx, profile)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (m *ProfileRepositoryMock) GetByUserID(ctx context.Context, userID int) (*models.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (m *ProfileRepositoryMock) GetByUserIDWithUser(ctx context.Context, userID int) (*models.Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Profile), args.Error(1)
}

func (m *ProfileRepositoryMock) Update(ctx context.Context, profile *models.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *ProfileRepositoryMock) UpdateByUserID(ctx context.Context, userID int, updates map[string]interface{}) error {
	args := m.Called(ctx, userID, updates)
	return args.Error(0)
}

func (m *ProfileRepositoryMock) DeleteByUserID(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type ProductRepositoryMock struct {
	mock.Mock
}

func (m *ProductRepositoryMock) Insert(ctx context.Context, product *models.Product) (*models.Product, error) {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetAll(ctx context.Context, page, limit int, search string, categoryID int) ([]models.Product, int64, error) {
	args := m.Called(ctx, page, limit, search, categoryID)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Product), args.Get(1).(int64), args.Error(2)
}

func (m *ProductRepositoryMock) Get(ctx context.Context, id int) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetBySlug(ctx context.Context, slug string) (*models.Product, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) Update(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *ProductRepositoryMock) GetByUser(ctx context.Context, userID int) ([]models.Product, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetByCategory(ctx context.Context, categoryID int) ([]models.Product, error) {
	args := m.Called(ctx, categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Product), args.Error(1)
}

type BasketRepositoryMock struct {
	mock.Mock
}

func (m *BasketRepositoryMock) GetActiveBasket(ctx context.Context, userID int) (*models.Basket, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Basket), args.Error(1)
}

func (m *BasketRepositoryMock) CreateBasket(ctx context.Context, basket *models.Basket) error {
	args := m.Called(ctx, basket)
	return args.Error(0)
}

func (m *BasketRepositoryMock) AddItem(ctx context.Context, basketID int, item *models.BasketItem) error {
	args := m.Called(ctx, basketID, item)
	return args.Error(0)
}

func (m *BasketRepositoryMock) UpdateItemQuantity(ctx context.Context, itemID int, quantity int) error {
	args := m.Called(ctx, itemID, quantity)
	return args.Error(0)
}

func (m *BasketRepositoryMock) RemoveItem(ctx context.Context, itemID int) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *BasketRepositoryMock) ClearBasket(ctx context.Context, basketID int) error {
	args := m.Called(ctx, basketID)
	return args.Error(0)
}
