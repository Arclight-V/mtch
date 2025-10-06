package models

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

type VerifyEmailResponse struct {
	User VerifiedEmailUserDTO
}
