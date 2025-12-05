package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" binding:"required,min=3" gorm:"uniqueIndex;not null"`
	Slug        string `json:"slug" gorm:"uniqueIndex"`
	Description string `json:"description" gorm:"type:text"`
	Image       string `json:"image"`
	Status      string `json:"status" gorm:"default:'active'"` // active, inactive

	ParentID *int       `json:"parentId"`
	Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`

	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}
