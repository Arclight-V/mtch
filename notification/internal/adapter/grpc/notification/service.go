package notification

import (
	usecase "github.com/Arclight-V/mtch/notification/internal/usecase/notification"
	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
)

type notificationServiceServer struct {
	notificationservicepb.UnimplementedNotificationServiceServer

	notificationUC usecase.NotificationUseCase
}

func NewNotificationServiceServer(notificationUC usecase.NotificationUseCase) *notificationServiceServer {
	return &notificationServiceServer{notificationUC: notificationUC}
}
