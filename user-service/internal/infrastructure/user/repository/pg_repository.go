package repository

import (
	"context"
	"log"
	"net/mail"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
)

// User Repository
type UserRepository struct {
	db *sqlx.DB
}

// NewUserRepository returns a new UserRepository
// TODO: replace db *FakeDB on *sqlx.DB
func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Create(ctx context.Context, regData *domain.RegisterInput) (*domain.User, error) {
	pendingUser, _ := domain.NewPendingUser(regData.PersonalDate)

	// TODO: Отказаться от StructScan
	if err := u.db.QueryRowxContext(
		ctx,
		createPendingUserQuery,
		pendingUser.PersonalData.FirstName,
		pendingUser.PersonalData.LastName,
		pendingUser.PersonalData.Contact,
		pendingUser.PersonalData.Email,
		pendingUser.PersonalData.Phone,
		pendingUser.PersonalData.Password,
		pendingUser.PersonalData.DateBirthday,
		pendingUser.PersonalData.Gender,
	).StructScan(pendingUser); err != nil {
		return nil, errors.Wrap(err, "Create.QueryRowxContext")
	}

	return pendingUser, nil

}

func (u *UserRepository) FindByContact(ctx context.Context, contact string) (*domain.User, error) {
	log.Println("FindByContact: ", contact)

	if _, err := mail.ParseAddress(contact); err != nil {
		return u.FindByEmail(ctx, contact)
	}
	return u.FindByEmail(ctx, contact)
}

func (u *UserRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	log.Println("FindByPhone: ", phone)

	panic("implement me")
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	log.Println("FindByEmail: ", email)

	user := &domain.User{PersonalData: &domain.PersonalData{
		Email: email,
	}}
	if err := u.db.GetContext(ctx, user, findByEmailQuery, email); err != nil {
		return nil, errors.Wrap(err, "FindByEmail.GetContext")
	}

	return user, nil
}
