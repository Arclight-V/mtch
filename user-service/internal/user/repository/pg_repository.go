package repository

import (
	"context"
	"github.com/pkg/errors"
	"user-service/internal/models"
)

// User Repository
type UserRepository struct {

	// TODO:: replace db *FakeDB on *sqlx.DB
	db *FakeDB
}

// NewUserRepository returns a new UserRepository
// TODO: replace db *FakeDB on *sqlx.DB
func NewUserRepository(db *FakeDB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) Create(ctx context.Context, regData *models.RegistrationData) (*models.User, error) {
	pendingUser, _ := models.NewPendingUser(regData.Email, regData.PasswordHash)

	if err := u.db.QueryRowxContext(
		ctx,
		createPendingUserQuery,
		pendingUser.Email,
		pendingUser.PasswordHash,
		pendingUser.UserID,
	).StructScan(regData); err != nil {
		return nil, errors.Wrap(err, "Create.QueryRowxContext")
	}
	return pendingUser, nil

}
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{Email: email}
	if err := u.db.GetContext(ctx, user, findByEmailQuery, email); err != nil {
		return nil, errors.Wrap(err, "FindByEmail.GetContext")
	}

	return user, nil
}
