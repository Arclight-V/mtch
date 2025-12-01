package repository

import (
	"context"
	"net/mail"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
)

// UserRepoDB is a wrapper around *sqlx.DB
type UserRepoDB struct {
	db *sqlx.DB

	logger log.Logger
}

// NewUserRepoDB returns a new UserRepoDB
func NewUserRepoDB(logger log.Logger, db *sqlx.DB) *UserRepoDB {
	logger = log.With(logger, "component", "UserRepoDB")

	return &UserRepoDB{
		db:     db,
		logger: logger,
	}
}

// Create creates new user in DB
func (u *UserRepoDB) Create(ctx context.Context, regData *domain.RegisterInput) (*domain.User, error) {
	level.Debug(u.logger).Log("msg", "create", "regData", regData)

	pendingUser, _ := domain.NewPendingUser(regData.PersonalDate)

	// TODO: Отказаться от StructScan
	if err := u.db.QueryRowxContext(
		ctx,
		createPendingUserQuery,
		pendingUser.UserID,
		pendingUser.FirstName,
		pendingUser.LastName,
		pendingUser.Contact,
		pendingUser.Phone,
		pendingUser.Email,
		pendingUser.Password,
		pendingUser.DateBirthday,
		pendingUser.Gender,
		pendingUser.Role,
	).StructScan(pendingUser); err != nil {
		return nil, errors.Wrap(err, "Create.QueryRowxContext")
	}

	return pendingUser, nil

}

// FindByContact finds user by contact
func (u *UserRepoDB) FindByContact(ctx context.Context, contact string) (*domain.User, error) {
	level.Debug(u.logger).Log("msg", "finding user by contact", "contact", contact)

	if _, err := mail.ParseAddress(contact); err != nil {
		return u.FindByEmail(ctx, contact)
	}
	return u.FindByEmail(ctx, contact)
}

// FindByPhone finds user by phone
func (u *UserRepoDB) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	level.Debug(u.logger).Log("msg", "finding user by phone", "phone", phone)

	panic("implement me")
}

// FindByEmail finds user by email
func (u *UserRepoDB) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	level.Debug(u.logger).Log("msg", "finding user by email", "email", email)

	user := &domain.User{PersonalData: domain.PersonalData{
		Email: email,
	}}
	if err := u.db.GetContext(ctx, user, findByEmailQuery, email); err != nil {
		return nil, errors.Wrap(err, "FindByEmail.GetContext")
	}

	return user, nil
}
