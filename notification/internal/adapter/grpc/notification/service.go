package notification

import (
	"github.com/go-kit/log"

	usecase "github.com/Arclight-V/mtch/notification/internal/usecase/notification"
	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
)

type notificationServiceServer struct {
	notificationservicepb.UnimplementedNotificationServiceServer

	notificationUC usecase.NotificationUseCase
	logger         log.Logger
}

func NewNotificationServiceServer(notificationUC usecase.NotificationUseCase, logger log.Logger) *notificationServiceServer {
	return &notificationServiceServer{notificationUC: notificationUC, logger: logger}
}
