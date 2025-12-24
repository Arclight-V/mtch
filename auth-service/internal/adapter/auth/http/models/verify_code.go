package models

import "time"

type VerifyCodeRequest struct {
	UserID string `json:"userid"`
	Code   string `json:"code"`
}

type VerifyCodeResponse struct {
	UserID     string    `json:"userid"`
	VerifiedAt time.Time `json:"verified_at"`
	Verified   bool      `json:"verified"`
}
