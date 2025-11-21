package notification

import (
	"github.com/go-kit/log"

	"github.com/Arclight-V/mtch/pkg/feature_list"

	usecase "github.com/Arclight-V/mtch/notification/internal/usecase/notification"
	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
)

type notificationServiceServer struct {
	notificationservicepb.UnimplementedNotificationServiceServer

	notificationUC usecase.NotificationUseCase
	logger         log.Logger
	featureList    *feature_list.FeatureList
}

func NewNotificationServiceServer(
	notificationUC usecase.NotificationUseCase,
	logger log.Logger,
	featureList *feature_list.FeatureList,
) *notificationServiceServer {
	return &notificationServiceServer{notificationUC: notificationUC, logger: logger, featureList: featureList}
}
