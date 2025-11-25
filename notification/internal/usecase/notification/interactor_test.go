package notification_test

import (
	"context"
	"testing"

	"github.com/go-kit/log"
	"github.com/golang/mock/gomock"

	"github.com/Arclight-V/mtch/pkg/feature_list"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
	"github.com/Arclight-V/mtch/notification/internal/features"
	"github.com/Arclight-V/mtch/notification/internal/infrastructure/codegen"
	notification "github.com/Arclight-V/mtch/notification/internal/usecase/notification"
	"github.com/Arclight-V/mtch/notification/internal/usecase/notification/mocks"
	"github.com/Arclight-V/mtch/notification/internal/usecase/repository"
)

func TestNotifyUserRegistered_Ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailSender := mocks.NewMockEmailSender(ctrl)
	logger := log.NewNopLogger()
	featureList := feature_list.NewNoopFeatureList(features.Features)
	verifyCodeMems := repository.NewVerifyCodesMem()
	codegen := codegen.NewNoopCodeGenerator()

	nuc := notification.NewNotificationUseCase(mockEmailSender, logger, featureList, verifyCodeMems, codegen)

	in := domain.Input{UserContacts: &domain.UserContacts{
		UserID: "u1",
		Contacts: []*domain.UserContact{
			{Channel: domain.ChannelEmail, Value: "a@b.com"},
			//{Channel: domain.ChanelPush, Value: "dev-token"},
		},
	}}

	vd := notification.VerifyData{
		Email:       in.UserContacts.Contacts[0].Value,
		VerifyToken: "token",
	}
	mockEmailSender.EXPECT().SendUserRegistered(context.Background(), vd).Return(nil).Times(1)

	out, err := nuc.NotifyUserRegistered(context.Background(), &in)

	if err != nil {
		t.Fatalf("Error sending email: %v", err)
	}
	if out == nil {
		t.Fatal("Output was nil")
	}
}
