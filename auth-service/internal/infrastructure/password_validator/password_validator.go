package password_validator

import (
	"errors"
	"github.com/go-passwd/validator"
)

type UserPasswordValidator struct {
	passwordValidator *validator.Validator
}

const (
	ErrPasswordTooShort       = "password must be at least 8 characters long"
	ErrPasswordTooLong        = "password must be no more than 20 characters long"
	ErrPasswordTooCommon      = "password is too common"
	ErrPasswordMissingUpper   = "password must contain at least one uppercase letter"
	ErrPasswordMissingLower   = "password must contain at least one lowercase letter"
	ErrPasswordMissingDigit   = "password must contain at least one digit"
	ErrPasswordMissingSpecial = "password must contain at least one special character"
)

func NewUserPasswordValidator() *UserPasswordValidator {
	return &UserPasswordValidator{
		passwordValidator: validator.New(
			validator.MinLength(8, errors.New(ErrPasswordTooShort)),
			validator.MaxLength(20, errors.New(ErrPasswordTooLong)),
			validator.CommonPassword(errors.New(ErrPasswordTooCommon)),
			validator.ContainsAtLeast("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 1, errors.New(ErrPasswordMissingUpper)),
			validator.ContainsAtLeast("abcdefghijklmnopqrstuvwxyz", 1, errors.New(ErrPasswordMissingLower)),
			validator.ContainsAtLeast("1234567890", 1, errors.New(ErrPasswordMissingDigit)),
			validator.ContainsAtLeast("/+-.:;\"'`<>,{}[]()^%$#@!~*", 1, errors.New(ErrPasswordMissingSpecial)),
		),
	}
}

func (u *UserPasswordValidator) Validate(password string) error {
	err := u.passwordValidator.Validate(password)
	if err != nil {
		return err
	}
	return nil
}
