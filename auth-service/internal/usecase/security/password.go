package security

//go:generate mockgen -source=$GOFILE -package=mocks -destination=../mocks/passwordHasher_mock.go
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hash, plain string) bool
}
