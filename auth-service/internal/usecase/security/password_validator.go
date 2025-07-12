package security

type PasswordValidator interface {
	Validate(password string) error
}
