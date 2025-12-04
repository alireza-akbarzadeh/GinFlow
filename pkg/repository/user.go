package repository

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        int        `json:"id" gorm:"primaryKey"`
	Email     string     `json:"email" gorm:"uniqueIndex;not null"`
	Name      string     `json:"name" gorm:"not null"`
	Password  string     `json:"-" gorm:"not null"` // Never expose password in JSON
	LastLogin *time.Time `json:"lastLogin"`
}

// UserRepository handles user database operations
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Insert creates a new user in the database
func (r *UserRepository) Insert(user *User) (*User, error) {
	result := r.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// Get retrieves a user by ID
func (r *UserRepository) Get(id int) (*User, error) {
	var user User
	result := r.DB.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*User, error) {
	var user User
	result := r.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetById retrieves a user by ID (alias for Get)
func (r *UserRepository) GetById(id int) (*User, error) {
	return r.Get(id)
}

// UpdatePassword updates the user's password
func (r *UserRepository) UpdatePassword(userID int, hashedPassword string) error {
	result := r.DB.Model(&User{}).Where("id = ?", userID).Update("password", hashedPassword)
	return result.Error
}

// GetAll retrieves all users
func (r *UserRepository) GetAll() ([]*User, error) {
	var users []*User
	result := r.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *User) error {
	result := r.DB.Save(user)
	return result.Error
}

// Delete removes a user by ID
func (r *UserRepository) Delete(id int) error {
	result := r.DB.Delete(&User{}, id)
	return result.Error
}

// UpdateLastLogin updates the last login timestamp for a user
func (r *UserRepository) UpdateLastLogin(id int) error {
	now := time.Now()
	result := r.DB.Model(&User{}).Where("id = ?", id).Update("last_login", now)
	return result.Error
}
