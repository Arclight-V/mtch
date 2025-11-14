package notification

import (
	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
)

type notificationServiceServer struct {
	notificationservicepb.UnimplementedNotificationServiceServer
}

func NewNotificationServiceServer() *notificationServiceServer {
	return &notificationServiceServer{}
}
