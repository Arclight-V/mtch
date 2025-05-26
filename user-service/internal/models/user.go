package models

import (
	"github.com/google/uuid"
	"time"
)

// User base model
type User struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id" validate:"omitempty"`
	Email     string    `json:"email" db:"email" validate:"omitempty,lte=60,email"`
	FirstName string    `json:"first_name" db:"first_name" validate:"required,lte=30"`
	LastName  string    `json:"last_name" db:"last_name" validate:"required,lte=30"`
	Role      string    `json:"role" db:"role" validate:"required"`
	// String?
	Avatar *string `json:"avatar" db:"avatar"`
	// String?
	PasswordHash string    `json:"password,omitempty" db:"password"`
	CreatedAt    time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" db:"updated_at"`
	Verified     bool      `json:"verified" db:"verified"`
}

// Create new pending User
func NewPendingUser(email, hash string) (*User, error) {
	// TODO
	// if !validator.IsEmail(email) {
	//	return nil, errors.New("invalid email")
	// }
	return &User{
		UserID:       uuid.New(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
		Verified:     false,
	}, nil
}

// Get avatar string
// TODO:: Use string?
func (u *User) GetAvatar() string {
	if u.Avatar == nil {
		return ""
	}
	return *u.Avatar
}
