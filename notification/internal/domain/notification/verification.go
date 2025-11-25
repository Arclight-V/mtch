package notification

import "time"

type VerificationPurpose int

const (
	EmailVerify VerificationPurpose = iota
)

func (verificationPurpose VerificationPurpose) String() string {
	switch verificationPurpose {
	case EmailVerify:
		return "email"

	default:
		return "unknown"
	}
}

// VerificationCode represents a one-time code used to confirm user actions.
type VerificationCode struct {
	UserID      string
	Code        string
	ExpiresAt   time.Time
	Purpose     VerificationPurpose
	Attempts    int
	MaxAttempts int
}
