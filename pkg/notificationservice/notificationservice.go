package notificationservice

import (
	"google.golang.org/grpc"

	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
)

func RegisterNotificationServer(notificationSrv notificationservicepb.NotificationServiceServer) func(*grpc.Server) {
	return func(s *grpc.Server) {
		notificationservicepb.RegisterNotificationServiceServer(s, notificationSrv)
	}
}
