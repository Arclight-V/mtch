package testdata

import (
	"time"

	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
)

func NewTestPendingUser() (*domain.User, error) {
	personalData := domain.PersonalData{
		FirstName:    "John",
		LastName:     "Doe",
		Contact:      "email",
		Phone:        "+7999999999",
		Email:        "a@b.com",
		Password:     "password",
		DateBirthday: time.Date(1992, time.Month(11), 28, 0, 0, 0, 0, time.UTC),
		Gender:       domain.Male,
	}

	pendingUser, err := domain.NewPendingUser(&personalData)

	return pendingUser, err
}
