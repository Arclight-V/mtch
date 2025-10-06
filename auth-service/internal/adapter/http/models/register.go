package models

type RegisterRequest struct {
	FirstName  string `json:"first_name" validate:"required"`
	LastName   string `json:"last_name" validate:"required"`
	Contact    string `json:"contact" validate:"required,contact"`
	Password   string `json:"password" validate:"required,min=8"`
	BirthDay   string `json:"birth_day" validate:"required"`
	BirthMonth string `json:"birth_month" validate:"required"`
	BirthYear  string `json:"birth_year" validate:"required"`
	Gender     string `json:"gender"`
}
type RegisterResponse struct {
	User PendingUserDTO
}
