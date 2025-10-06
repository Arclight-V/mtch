package auth

import (
	"context"
	"github.com/Arclight-V/mtch/auth-service/internal/usecase/security"
	"net/mail"
	"regexp"
)

type Date struct {
	BirthDay   int32
	BirthMonth int32
	BirthYear  int32
}
type RegisterInput struct {
	FirstName string
	LastName  string
	Contact   string
	Phone     string
	Email     string
	Password  string
	Date      *Date
	Gender    string
}

var phoneRegex = regexp.MustCompile(`^\+?[0-9\s\-\(\)]{7,20}$`)

func (ri *RegisterInput) SetPassword(plain string, h security.PasswordHasher) error {
	hash, err := h.Hash(plain)
	if err != nil {
		return err
	}
	ri.Password = hash
	return nil
}

func (ri *RegisterInput) SetEmailOrPhone() error {
	if _, err := mail.ParseAddress(ri.Contact); err != nil {
		if phoneRegex.MatchString(ri.Contact) {
			ri.Phone = ri.Contact
			return nil
		} else {
			return err
		}
	}
	ri.Email = ri.Contact

	return nil
}

type RegisterOutput struct {
	UserID   string
	Email    string
	Verified bool
}

//go:generate mockgen -source=$GOFILE -package=mocks -destination=../mocks/register_mock.go
type RegisterUseCase interface {
	Register(ctx context.Context, in *RegisterInput) (*RegisterOutput, error)
}
