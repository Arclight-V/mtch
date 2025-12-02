package user

import (
	"github.com/go-kit/log"

	"github.com/Arclight-V/mtch/pkg/feature_list"
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"

	usecase "github.com/Arclight-V/mtch/user-service/internal/usecase/user"
)

type usersServiceServer struct {
	userservicepb.UnimplementedUserServiceServer

	userUC      usecase.UserUseCase
	logger      log.Logger
	featureList *feature_list.FeatureList
}

func NewUserServiceServer(
	userUC usecase.UserUseCase,
	logger log.Logger,
	featureList *feature_list.FeatureList,
) *usersServiceServer {
	logger = log.With(logger, "component", "gRPC/usersServiceServer")

	return &usersServiceServer{
		userUC:      userUC,
		logger:      logger,
		featureList: featureList,
	}
}
