package models

// Category represents an event category
type Category struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" binding:"required,min=3" gorm:"uniqueIndex;not null"`
	Description string `json:"description"`
}
