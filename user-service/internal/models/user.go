package models

import (
	"encoding/json"
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

// NewPendingUser Create new pending User
func NewPendingUser(email, hash string) (*User, error) {
	return &User{
		UserID:       uuid.New(),
		Email:        email,
		PasswordHash: hash,
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

// RegistrationData data for registration
type RegistrationData struct {
	Email        string
	PasswordHash string
}

type CreateUserStatus int

const (
	CreateUserStatusUnspecified = iota

	//CreatedUnverified Successfully created, but not yet verified
	CreatedUnverified

	//ExistsVerified There is already a user with this email address and it has been verified
	ExistsVerified

	//ExistsUnverified Already exists, but has NOT been verified
	ExistsUnverified

	//Rejected Not created for a business reason
	Rejected
)

func (c CreateUserStatus) String() string {
	switch c {
	case CreateUserStatusUnspecified:
		return "CREATE_USER_STATUS_UNSPECIFIED"
	case CreatedUnverified:
		return "CREATED_UNVERIFIED"
	case ExistsVerified:
		return "EXISTS_VERIFIED"
	case ExistsUnverified:
		return "EXISTS_UNVERIFIED"
	case Rejected:
		return "REJECTED"
	}
	return ""
}

func (c CreateUserStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

type RegistrationOutput struct {
	UserID uuid.UUID
	Status CreateUserStatus
}
