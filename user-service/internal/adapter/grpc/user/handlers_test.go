package user

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/go-kit/log"
	"github.com/golang/mock/gomock"

	"github.com/Arclight-V/mtch/pkg/feature_list"
	"github.com/Arclight-V/mtch/pkg/userservice/userservicepb/v1"

	"github.com/Arclight-V/mtch/user-service/internal/adapter/grpc/user/testdata"
	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
	"github.com/Arclight-V/mtch/user-service/internal/usecase/user"
	"github.com/Arclight-V/mtch/user-service/internal/usecase/user/mocks"
)

func TestUsersServiceServer_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserUseCase := mocks.NewMockUserUseCase(ctrl)
	logger := log.NewNopLogger()
	featureList := feature_list.NewNoopFeatureList(feature_list.Features{})

	nuss := NewUserServiceServer(mockUserUseCase, logger, featureList)

	req := testdata.NewTestPBRequest()

	tests := []struct {
		name      string
		req       *userservicepb.RegisterRequest
		regOutput *domain.RegisterOutput
		wantErr   bool
		err       error
	}{
		{
			name:      "should create a new user",
			req:       req,
			regOutput: &domain.RegisterOutput{},
			wantErr:   false,
			err:       nil,
		},
		{
			name:      "should fail to create a user 1",
			req:       req,
			regOutput: nil,
			wantErr:   true,
			err:       user.ErrUserIsExist,
		},
		{
			name:      "should fail to create a user 2",
			req:       req,
			regOutput: nil,
			wantErr:   true,
			err:       user.ErrUserIsExistUnverified,
		},
		{
			name:      "should fail to create a user 3",
			req:       req,
			regOutput: nil,
			wantErr:   true,
			err:       user.ErrUserNotCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserUseCase.EXPECT().
				Register(gomock.Any(), gomock.AssignableToTypeOf(&domain.RegisterInput{})).
				Return(tt.regOutput, tt.err).
				Times(1)

			resp, err := nuss.Register(context.Background(), tt.req)
			if tt.wantErr {
				assert.NotNil(t, err)
				if !errors.Is(err, tt.err) {
					t.Errorf("got %v; want %v", err, tt.err)
				}
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
			}
		})
	}
}
