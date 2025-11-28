package notification

import (
	"context"
	"github.com/Arclight-V/mtch/notification/internal/infrastructure/repository"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"

	"github.com/Arclight-V/mtch/pkg/feature_list"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
	"github.com/Arclight-V/mtch/notification/internal/features"
	"github.com/Arclight-V/mtch/notification/internal/infrastructure/codegen"
)

type notificationUseCase struct {
	emailSender   EmailSender
	varifyCodeMem *repository.VerifyCodesMem
	repo          PGRepository

	logger        log.Logger
	featureList   *feature_list.FeatureList
	codeGenerator *codegen.CodeGenerator
}

func NewNotificationUseCase(
	emailSender EmailSender,
	logger log.Logger,
	featureList *feature_list.FeatureList,
	verifyCodeMem *repository.VerifyCodesMem,
	codeGenerator *codegen.CodeGenerator,
	repo PGRepository,
) *notificationUseCase {
	return &notificationUseCase{
		emailSender:   emailSender,
		logger:        logger,
		featureList:   featureList,
		varifyCodeMem: verifyCodeMem,
		codeGenerator: codeGenerator,
		repo:          repo,
	}
}

func (n *notificationUseCase) NotifyUserRegistered(ctx context.Context, in *domain.Input) (*domain.Output, error) {
	level.Info(n.logger).Log("msg", "NotifyUserRegistered:", "domain.Input", in)

	var err error
	vc := n.codeGenerator.NewVerificationCode(in.UserContacts.UserID)

	for _, c := range in.UserContacts.Contacts {
		switch c.Channel {
		case domain.ChannelEmail:
			if !n.featureList.IsEnabled(features.StoreCodesInDB) {
				if errInsert := n.varifyCodeMem.InsertIssue(ctx, vc); errInsert != nil {
					err = errors.Wrap(err, errInsert.Error())
				}
			} else {
				if errInsert := n.repo.InsertIssue(ctx, vc); errInsert != nil {
					err = errors.Wrap(err, errInsert.Error())
				}
			}
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
