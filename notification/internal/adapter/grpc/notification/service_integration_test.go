package notification

import (
	"context"
	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
	"github.com/Arclight-V/mtch/notification/internal/usecase/notification/mocks"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"net"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/golang/mock/gomock"

	"google.golang.org/grpc"

	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
)

const bufSize = 1024 * 1024

func newBufconnServer(t *testing.T) (context.Context, *grpc.ClientConn, func()) {
	lis := bufconn.Listen(bufSize)

	s := grpc.NewServer()
	mockCtrl := gomock.NewController(t)

	mockNotificationUseCase := mocks.NewMockNotificationUseCase(mockCtrl)
	req := notificationservicepb.NotificationUserContactsRequest{
		UserID: "u1",
		Contacts: []*notificationservicepb.Contact{
			{Chanel: notificationservicepb.Channel_ChannelEmail, Value: "a@b.com"},
			{Chanel: notificationservicepb.Channel_ChannelPush, Value: "dev-token"},
		},
	}
	in := &domain.Input{UserContacts: protoContactsToUserContacts(req.UserID, req.Contacts)}
	mockNotificationUseCase.EXPECT().NotifyUserRegistered(gomock.Any(), in).
		Return(&domain.Output{}, nil).Times(1)

	logger := log.NewNopLogger()
	nss := NewNotificationServiceServer(mockNotificationUseCase, logger)

	notificationservicepb.RegisterNotificationServiceServer(s, nss)

	go func() {
		_ = s.Serve(lis)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	dialer := func(ctx context.Context, _ string) (net.Conn, error) {
		return lis.DialContext(ctx)
	}

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		mockCtrl.Finish()
		cancel()
		_ = conn.Close()
		s.Stop()
		_ = lis.Close()
	}

	return ctx, conn, cleanup
}

func TestNotificationService_Integration(t *testing.T) {
	ctx, conn, cleanup := newBufconnServer(t)
	defer cleanup()

	client := notificationservicepb.NewNotificationServiceClient(conn)

	req := notificationservicepb.NotificationUserContactsRequest{
		UserID: "u1",
		Contacts: []*notificationservicepb.Contact{
			{Chanel: notificationservicepb.Channel_ChannelEmail, Value: "a@b.com"},
			{Chanel: notificationservicepb.Channel_ChannelPush, Value: "dev-token"},
		},
	}

	resp, err := client.NotifyUserRegistered(ctx, &req)

	if err != nil {
		t.Fatalf("NotificationService.NotifyUserRegistered(): %v", err)
	}

	if resp == nil {
		t.Fatalf("NotificationService.NotifyUserRegistered(): nil response")
	}

}
