package models

import (
	"time"
)

type PendingUserDTO struct {
	UserID   string    `json:"userid"`
	Email    string    `json:"email"`
	CreateAt time.Time `json:"create_at"`
	Verified bool      `json:"verified"`
}

type VerifiedEmailUserDTO struct {
	UserID     string    `json:"userid"`
	VerifiedAt time.Time `json:"activate_at"`
	Verified   bool      `json:"verified"`
}
