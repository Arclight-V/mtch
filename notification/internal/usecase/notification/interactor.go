package notification

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
)

type notificationUseCase struct {
	logger log.Logger
}

func NewNotificationUseCase(logger log.Logger) *notificationUseCase {
	return &notificationUseCase{logger: logger}
}

func (u *notificationUseCase) NoopMethod(ctx context.Context, in *domain.NoopInput) (*domain.NoopOutput, error) {
	level.Info(u.logger).Log("msg", "NoopMethod called:", "NoopString", in.NoopStruct.NoopString)

	return &domain.NoopOutput{NoopString: "NoopString"}, nil
}
