package notification

import (
	"context"
	"log"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
)

type notificationUseCase struct {
}

func NewNotificationUseCase() *notificationUseCase {
	return &notificationUseCase{}
}

func (u *notificationUseCase) NoopMethod(ctx context.Context, in *domain.NoopInput) (*domain.NoopOutput, error) {
	log.Println("NoopMethod called:", in.NoopStruct.NoopString)

	return &domain.NoopOutput{NoopString: "NoopString"}, nil
}
