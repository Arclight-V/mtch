package repository

import (
	"context"
	"errors"
	"net/mail"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
)

// UsersDBMem temporary user storage in memory of database unavailability
type UsersDBMem struct {
	mu    sync.Mutex
	users map[string]domain.User

	logger log.Logger
}

// NewUsersDBMem returns new *UsersDBMem
func NewUsersDBMem(logger log.Logger) *UsersDBMem {
	logger = log.With(logger, "component", "UsersDBMem")

	return &UsersDBMem{
		users:  make(map[string]domain.User),
		logger: logger,
	}
}

// Create creates create new user in db
func (m *UsersDBMem) Create(ctx context.Context, regData *domain.RegisterInput) (*domain.User, error) {
	level.Debug(m.logger).Log("msg", "create", "regData", regData)

	pendingUser, _ := domain.NewPendingUser(regData.PersonalDate)
	m.mu.Lock()
	defer m.mu.Unlock()

	cnt := pendingUser.PersonalData.Contact
	if _, ok := m.users[cnt]; ok {
		return nil, errors.New("userservice with " + cnt + " already exists")
	}
	m.users[cnt] = *pendingUser

	return pendingUser, nil
}

// FindByContact finds user by contact
func (m *UsersDBMem) FindByContact(ctx context.Context, contact string) (*domain.User, error) {
	level.Debug(m.logger).Log("msg", "finding user by contact", "contact", contact)

	if _, err := mail.ParseAddress(contact); err != nil {
		return m.FindByPhone(ctx, contact)
	}
	return m.FindByEmail(ctx, contact)
}

// FindByEmail finds user by email
func (m *UsersDBMem) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	level.Debug(m.logger).Log("msg", "finding user by email", "email", email)

	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[email]
	if !ok {
		return nil, errors.New("userservice not found")
	}
	return &user, nil
}

// FindByPhone finds users by phone
func (m *UsersDBMem) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	level.Debug(m.logger).Log("msg", "finding user by phone", "phone", phone)

	m.mu.Lock()
	defer m.mu.Unlock()

	user, ok := m.users[phone]
	if !ok {
		return nil, errors.New("userservice not found")
	}
	return &user, nil
}
