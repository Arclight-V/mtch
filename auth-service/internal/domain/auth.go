package domain

import "time"

type VerifyTokenIssue struct {
	JTI       string
	UserID    string
	ExpiresAt time.Time
}
