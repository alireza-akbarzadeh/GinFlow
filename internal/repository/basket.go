package repository

import (
	"errors"

	"github.com/alireza-akbarzadeh/ginflow/internal/models"
	"gorm.io/gorm"
)

type BasketRepository struct {
	DB *gorm.DB
}

func NewBasketRepository(db *gorm.DB) *BasketRepository {
	return &BasketRepository{DB: db}
}

// GetActiveBasket retrieves the active basket for a user
func (r *BasketRepository) GetActiveBasket(userID int) (*models.Basket, error) {
	var basket models.Basket
	err := r.DB.Preload("Items.Product").Where("user_id = ? AND status = ?", userID, "active").First(&basket).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &basket, nil
}

// CreateBasket creates a new basket
func (r *BasketRepository) CreateBasket(basket *models.Basket) error {
	return r.DB.Create(basket).Error
}

// AddItem adds an item to the basket or updates quantity if it exists
func (r *BasketRepository) AddItem(basketID int, item *models.BasketItem) error {
	var existingItem models.BasketItem
	err := r.DB.Where("basket_id = ? AND product_id = ?", basketID, item.ProductID).First(&existingItem).Error

	if err == nil {
		// Item exists, update quantity
		existingItem.Quantity += item.Quantity
		return r.DB.Save(&existingItem).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// Item does not exist, create new
		item.BasketID = basketID
		return r.DB.Create(item).Error
	}
	return err
}

// UpdateItemQuantity updates the quantity of an item in the basket
func (r *BasketRepository) UpdateItemQuantity(itemID int, quantity int) error {
	if quantity <= 0 {
		return r.DB.Delete(&models.BasketItem{}, itemID).Error
	}
	return r.DB.Model(&models.BasketItem{}).Where("id = ?", itemID).Update("quantity", quantity).Error
}

// RemoveItem removes an item from the basket
func (r *BasketRepository) RemoveItem(itemID int) error {
	return r.DB.Delete(&models.BasketItem{}, itemID).Error
}

// ClearBasket removes all items from the basket
func (r *BasketRepository) ClearBasket(basketID int) error {
	return r.DB.Where("basket_id = ?", basketID).Delete(&models.BasketItem{}).Error
}
