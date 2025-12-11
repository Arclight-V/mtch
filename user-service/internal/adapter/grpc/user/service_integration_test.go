package user

import (
	"context"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/Arclight-V/mtch/pkg/feature_list"
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"

	"github.com/Arclight-V/mtch/user-service/internal/adapter/grpc/user/testdata"
	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
	"github.com/Arclight-V/mtch/user-service/internal/features"
	"github.com/Arclight-V/mtch/user-service/internal/usecase/user/mocks"
)

const bufSize = 1024 * 1024

func newBufConnServer(t *testing.T) (context.Context, *grpc.ClientConn, func()) {
	lis := bufconn.Listen(bufSize)

	s := grpc.NewServer()
	mockCtrl := gomock.NewController(t)

	mockUserUseCase := mocks.NewMockUserUseCase(mockCtrl)

	req := testdata.NewTestPBRequest()
	pd := domain.NewPersonalDataFromRegisterRequest(req)
	in := &domain.RegisterInput{PersonalDate: pd}
	mockUserUseCase.EXPECT().Register(gomock.Any(), in).
		Return(&domain.RegisterOutput{}, nil).Times(1)

	logger := log.NewNopLogger()
	featureList := feature_list.NewNoopFeatureList(features.Features)

	nuss := NewUserServiceServer(mockUserUseCase, logger, featureList)

	userservicepb.RegisterUserServiceServer(s, nuss)

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

func TestUserService_Register_Integration(t *testing.T) {
	ctx, conn, cleanup := newBufConnServer(t)
	defer cleanup()

	client := userservicepb.NewUserServiceClient(conn)

	req := testdata.NewTestPBRequest()
	resp, err := client.Register(ctx, req)
	if err != nil {
		t.Fatalf("UserService.Register failed: %v", err)
	}
	if resp == nil {
		t.Fatalf("UserService.Register returned nil response")
	}
}
