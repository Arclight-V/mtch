package repository

import (
	"context"
	"errors"
	"log"
	"net/mail"
	"sync"

	domain "user-service/internal/domain/user"
)

type UsersDBMemory struct {
	mu    sync.Mutex
	users map[string]domain.User
}

func NewUsersDBMemory() *UsersDBMemory {
	return &UsersDBMemory{users: make(map[string]domain.User)}
}

func (m *UsersDBMemory) Create(ctx context.Context, regData *domain.RegisterInput) (*domain.User, error) {
	log.Println("Create: ", regData)

	pendingUser, _ := domain.NewPendingUser(regData.PersonalDate)
	m.mu.Lock()
	defer m.mu.Unlock()

	cnt := pendingUser.PersonalData.Contact
	if _, ok := m.users[cnt]; ok {
		return nil, errors.New("user with " + cnt + " already exists")
	}
	m.users[cnt] = *pendingUser

	return pendingUser, nil
}

func (m *UsersDBMemory) FindByContact(ctx context.Context, contact string) (*domain.User, error) {
	log.Println("FindByContact: ", contact)

	if _, err := mail.ParseAddress(contact); err != nil {
		return m.FindByPhone(ctx, contact)
	}
	return m.FindByEmail(ctx, contact)
}

func (m *UsersDBMemory) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	log.Println("FindByEmail: ", email)

	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[email]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (m *UsersDBMemory) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	log.Println("FindByPhone: ", phone)

	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[phone]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
