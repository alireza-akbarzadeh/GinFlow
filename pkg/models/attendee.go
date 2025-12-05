package models

// Attendee represents an attendee relationship between a user and an event
type Attendee struct {
	ID      int `json:"id" gorm:"primaryKey"`
	UserID  int `json:"userId" gorm:"not null"`
	EventID int `json:"eventId" gorm:"not null"`
}
