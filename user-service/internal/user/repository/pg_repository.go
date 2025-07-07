package repository

import (
	"context"
	"log"
	"user-service/internal/models"
)

// User Repository
type UserRepository struct {
	// TODO:: Add
	// db *sqlx.DB?
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u *UserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	// TODO:: Add logic

	return user, nil
}
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	// TODO:: Add logic
	log.Println("FindByEmail")
	user := &models.User{Email: email}

	return user, nil
}
