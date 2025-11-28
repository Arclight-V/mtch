package repository

import (
	"context"
	"errors"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
)

type VerifyCodesMem struct {
	mu     sync.Mutex
	byCode map[string][]domain.VerificationCode

	logger log.Logger
}

func NewVerifyCodesMem(logger log.Logger) *VerifyCodesMem {
	logger = log.With(logger, "component", "VerifyCodesMem")

	return &VerifyCodesMem{
		logger: logger,
		byCode: make(map[string][]domain.VerificationCode),
	}
}

func (m *VerifyCodesMem) InsertIssue(_ context.Context, v *domain.VerificationCode) error {
	level.Debug(m.logger).Log("msg", "inserting issue", "v", v)

	m.mu.Lock()
	defer m.mu.Unlock()

	if codes, ok := m.byCode[v.UserID]; ok {
		for i := range codes {
			if codes[i].Code == v.Code {
				return errors.New("verify code already exists")
			}
		}
	}
	m.byCode[v.UserID] = append(m.byCode[v.UserID], *v)

	return nil
}
