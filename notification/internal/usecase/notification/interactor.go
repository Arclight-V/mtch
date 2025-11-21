package notification

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"

	"github.com/Arclight-V/mtch/pkg/feature_list"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
)

type notificationUseCase struct {
	emailSender EmailSender

	logger      log.Logger
	featureList *feature_list.FeatureList
}

func NewNotificationUseCase(
	emailSender EmailSender,
	logger log.Logger,
	featureList *feature_list.FeatureList) *notificationUseCase {
	return &notificationUseCase{emailSender: emailSender, logger: logger, featureList: featureList}
}

func (n *notificationUseCase) NotifyUserRegistered(ctx context.Context, in *domain.Input) (*domain.Output, error) {
	level.Info(n.logger).Log("msg", "NotifyUserRegistered:", "domain.Input", in)

	var err error

	for _, c := range in.UserContacts.Contacts {
		switch c.Channel {
		case domain.ChannelEmail:
			vd := VerifyData{
				Email:       c.Value,
				VerifyToken: "token",
			}
			if sendErr := n.emailSender.SendUserRegistered(ctx, vd); sendErr != nil {
				_ = errors.Wrap(err, sendErr.Error())
			}

		case domain.ChanelPush:
		case domain.ChannelCall:
		case domain.Reject:
		}
	}

	if err != nil {
		return nil, err
	}
	return &domain.Output{}, nil
}
