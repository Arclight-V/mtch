package notification

import (
	"context"

	"github.com/go-kit/log/level"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
)

func (s *notificationServiceServer) NotifyUserRegistered(ctx context.Context, req *notificationservicepb.NotificationUserContactsRequest) (*notificationservicepb.NotificationUserContactsResponse, error) {
	level.Info(s.logger).Log("msg", "NoopMethod called:", "NoopString", req)

	usc := protoContactsToUserContacts(req.UserID, req.Contacts)
	in := &domain.Input{UserContacts: usc}

	_, err := s.notificationUC.NotifyUserRegistered(ctx, in)
	if err != nil {
		return &notificationservicepb.NotificationUserContactsResponse{}, err
	}
	resp := &notificationservicepb.NotificationUserContactsResponse{}

	return resp, nil
}
