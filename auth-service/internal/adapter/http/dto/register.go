package dto

type RegisterRequest struct {
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
}
type RegisterResponse struct {
	User PendingUserDTO
}
