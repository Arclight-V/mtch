package grpc

import (
	"github.com/go-kit/log"

	"github.com/Arclight-V/mtch/pkg/feature_list"

	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"
)

type AuthClient struct {
	userClient         userservicepb.UserServiceClient
	notificationClient notificationservicepb.NotificationServiceClient

	logger      log.Logger
	featureList *feature_list.FeatureList
}

func NewAuthClient(
	userClient userservicepb.UserServiceClient,
	notificationClient notificationservicepb.NotificationServiceClient,
	logger log.Logger,
	featureList *feature_list.FeatureList,
) *AuthClient {
	logger = log.With(logger, "component", "gRPC/GrpcUserClient")

	return &AuthClient{
		userClient:         userClient,
		notificationClient: notificationClient,
		logger:             logger,
		featureList:        featureList,
	}
}
