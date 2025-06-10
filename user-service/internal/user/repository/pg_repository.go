package repository

import (
	"context"
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
