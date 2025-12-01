package user

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type CreateUserStatus int

const (
	CreateUserStatusUnspecified = iota

	//CreatedUnverified Successfully created, but not yet verified
	CreatedUnverified

	//ExistsVerified There is already a userservice with this email address and it has been verified
	ExistsVerified

	//ExistsUnverified Already exists, but has NOT been verified
	ExistsUnverified

	//Rejected Not created for a business reason
	Rejected
)

type Gender int

const (
	Male Gender = iota
	Female
)

// User user data
type User struct {
	PersonalData
	UserID    uuid.UUID `db:"user_id"`
	Role      string    `db:"role"`
	Avatar    *string   `db:"avatar"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Activated bool      `db:"activated"`
}

// PersonalData personal user data
type PersonalData struct {
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
	Contact      string    `db:"contact"`
	Phone        string    `db:"phone"`
	Email        string    `db:"email"`
	Password     string    `db:"password"`
	DateBirthday time.Time `db:"date_birthday"`
	Gender       Gender    `db:"gender"`
}

// NewPendingUser Create new pending User
func NewPendingUser(data *PersonalData) (*User, error) {
	return &User{
		PersonalData: *data,
		Role:         "pending",
		Activated:    false,
		//TODO: move to db
		UserID: uuid.New(),
	}, nil
}

func (p *PersonalData) SetDateBirthday(year, month, day int) {
	p.DateBirthday = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// Get avatar string
// TODO:: Use string?
func (u *User) GetAvatar() string {
	if u.Avatar == nil {
		return ""
	}
	return *u.Avatar
}

// RegisterInput
type RegisterInput struct {
	PersonalDate *PersonalData
}

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

type RegisterOutput struct {
	UserID uuid.UUID
	Status CreateUserStatus
}

type VerifyEmailInput struct {
	UserID string
}

type VerifyEmailOutput struct {
	VerifiedAt time.Time
	Verified   bool
}
