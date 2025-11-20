package notification

import (
	"context"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
)

// NotificationUseCase interface
//
//go:generate mockgen -source=$GOFILE -package=mocks -destination=./mocks/ports_mock.go
type NotificationUseCase interface {
	NotifyUserRegistered(ctx context.Context, in *domain.Input) (*domain.Output, error)
}
