package notification_test

import (
	"context"
	"testing"

	"github.com/go-kit/log"
	"github.com/golang/mock/gomock"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
	notification "github.com/Arclight-V/mtch/notification/internal/usecase/notification"
	"github.com/Arclight-V/mtch/notification/internal/usecase/notification/mocks"
)

func TestNotifyUserRegistered_Ok(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailSender := mocks.NewMockEmailSender(ctrl)
	logger := log.NewNopLogger()

	nuc := notification.NewNotificationUseCase(mockEmailSender, logger)

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
