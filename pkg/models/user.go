package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        int        `json:"id" gorm:"primaryKey"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null"`
	Name      string     `json:"name" gorm:"not null"`
	Password  string     `json:"-" gorm:"not null"` // Never expose password in JSON
	LastLogin *time.Time `json:"lastLogin"`
}
