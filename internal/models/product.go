package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Product struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name" binding:"required,min=3" gorm:"not null"`
	Description string  `json:"description" gorm:"type:text"`
	Price       float64 `json:"price" binding:"required,gt=0" gorm:"not null"`
	Stock       int     `json:"stock" binding:"required,gte=0" gorm:"not null"`

	// Advanced Product Details
	SKU    string         `json:"sku" gorm:"unique;not null"`
	Status string         `json:"status" gorm:"default:'active'"`
	Slug   string         `json:"slug" gorm:"uniqueIndex"`
	Image  string         `json:"image"`
	Images pq.StringArray `json:"images" gorm:"type:text[]" swaggertype:"array,string" example:"[\"url1\",\"url2\"]"`
	Tags   pq.StringArray `json:"tags" gorm:"type:text[]" swaggertype:"array,string" example:"[\"tag1\",\"tag2\"]"`

	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`

	Discount   float64 `json:"discount" gorm:"default:0"`
	FinalPrice float64 `json:"finalPrice"`
	Brand      string  `json:"brand"`
	Weight     float64 `json:"weight"`
	Dimensions string  `json:"dimensions"`

	Rating       float64 `json:"rating" gorm:"default:0"`
	ReviewsCount int     `json:"reviewsCount" gorm:"default:0"`
	Views        int     `json:"views" gorm:"default:0"`

	UserID int  `json:"userId" gorm:"not null"`
	User   User `json:"user" gorm:"foreignKey:UserID"`

	Categories []Category `json:"categories" gorm:"many2many:product_categories;"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty" swaggerignore:"true"`
}
