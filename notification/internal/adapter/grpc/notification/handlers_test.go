package notification

import (
	"context"
	"testing"

	"github.com/go-kit/log"
	"github.com/golang/mock/gomock"

	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
	"github.com/Arclight-V/mtch/notification/internal/usecase/notification/mocks"
)

func TestNotifyUserRegistered_OK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockNotificationUseCase := mocks.NewMockNotificationUseCase(mockCtrl)
	logger := log.NewNopLogger()
	nss := NewNotificationServiceServer(mockNotificationUseCase, logger)

	req := notificationservicepb.NotificationUserContactsRequest{
		UserID: "u1",
		Contacts: []*notificationservicepb.Contact{
			{Chanel: notificationservicepb.Channel_ChannelEmail, Value: "a@b.com"},
			{Chanel: notificationservicepb.Channel_ChannelPush, Value: "dev-token"},
		},
	}

	in := &domain.Input{UserContacts: protoContactsToUserContacts(req.UserID, req.Contacts)}
	mockNotificationUseCase.EXPECT().NotifyUserRegistered(context.Background(), in).
		Return(&domain.Output{}, nil).Times(1)

	resp, err := nss.NotifyUserRegistered(context.Background(), &req)
	if err != nil {
		t.Fatalf("NotifyUserRegistered failed: %v", err)
	}
	if resp == nil {
		t.Fatalf("NotifyUserRegistered should have returned nil response")
	}

}
