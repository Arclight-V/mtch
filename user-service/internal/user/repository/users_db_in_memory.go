package repository

import (
	"context"
	"errors"
	"log"
	"net/mail"
	"sync"
	"user-service/internal/models"
)

type UsersDBMemory struct {
	mu    sync.Mutex
	users map[string]models.User
}

func NewUsersDBMemory() *UsersDBMemory {
	return &UsersDBMemory{users: make(map[string]models.User)}
}

func (m *UsersDBMemory) Create(ctx context.Context, regData *models.RegisterInput) (*models.User, error) {
	log.Println("Create: ", regData)

	pendingUser, _ := models.NewPendingUser(regData.PersonalDate)
	m.mu.Lock()
	defer m.mu.Unlock()

	cnt := pendingUser.PersonalData.Contact
	if _, ok := m.users[cnt]; ok {
		return nil, errors.New("user with " + cnt + " already exists")
	}
	m.users[cnt] = *pendingUser

	return pendingUser, nil
}

func (m *UsersDBMemory) FindByContact(ctx context.Context, contact string) (*models.User, error) {
	log.Println("FindByContact: ", contact)

	if _, err := mail.ParseAddress(contact); err != nil {
		return m.FindByPhone(ctx, contact)
	}
	return m.FindByEmail(ctx, contact)
}

func (m *UsersDBMemory) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	log.Println("FindByEmail: ", email)

	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (m *UsersDBMemory) FindByPhone(ctx context.Context, phone string) (*models.User, error) {
	log.Println("FindByPhone: ", phone)

	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[phone]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
