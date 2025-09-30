package dto

import (
	"time"
)

type PendingUserDTO struct {
	UserID   string    `json:"userid"`
	Email    string    `json:"email"`
	CreateAt time.Time `json:"create_at"`
	Verified bool      `json:"verified"`
}
