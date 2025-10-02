package repository

import (
	"context"
	"errors"
	"github.com/Arclight-V/mtch/auth-service/internal/domain"
	"log"
	"sync"
)

type VerifyTokensMem struct {
	mu    sync.Mutex
	byJTI map[string]domain.VerifyTokenIssue
}

func NewVerifyTokensMem() *VerifyTokensMem {
	return &VerifyTokensMem{byJTI: make(map[string]domain.VerifyTokenIssue)}
}

func (m *VerifyTokensMem) InsertIssue(_ context.Context, v domain.VerifyTokenIssue) error {
	log.Println("InsertIssue: ", v)
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.byJTI[v.JTI]; ok {
		return errors.New("verify token jti already exists")
	}
	m.byJTI[v.JTI] = v

	return nil
}

func (m *VerifyTokensMem) TryConsumeJTI(ctx context.Context, v domain.VerifyEmailToken) error {
	log.Println("TryConsumeJTI: ", v)
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.byJTI[v.JTI]; ok {
		delete(m.byJTI, v.JTI)
	} else {
		return errors.New("verify token jti not found")
	}

	return nil
}
