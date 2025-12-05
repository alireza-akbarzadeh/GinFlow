package models

import (
	"time"

	"gorm.io/gorm"
)

type Basket struct {
	ID        int            `json:"id" gorm:"primaryKey"`
	UserID    *int           `json:"userId" gorm:"index"`
	User      *User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Status    string         `json:"status" gorm:"default:'active'"` // active, completed
	Items     []BasketItem   `json:"items" gorm:"foreignKey:BasketID"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty" swaggerignore:"true"`
}

type BasketItem struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	BasketID  int       `json:"basketId" gorm:"not null;index"`
	ProductID int       `json:"productId" gorm:"not null"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int       `json:"quantity" gorm:"not null;check:quantity > 0"`
	UnitPrice float64   `json:"unitPrice" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
