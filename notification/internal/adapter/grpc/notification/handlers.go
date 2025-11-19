package notification

import (
	"context"

	"github.com/go-kit/log/level"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
	"github.com/Arclight-V/mtch/pkg/notificationservice/notificationservicepb/v1"
)

func (s *notificationServiceServer) NoopMethod(ctx context.Context, req *notificationservicepb.NoopRequest) (*notificationservicepb.NoopResponse, error) {
	level.Info(s.logger).Log("msg", "NoopMethod called:", "NoopString", req.NoopString)

	ns := &domain.NoopStruct{
		NoopString: req.NoopString,
	}

	in := &domain.NoopInput{NoopStruct: ns}

	noop, err := s.notificationUC.NoopMethod(ctx, in)
	if err != nil {
		return &notificationservicepb.NoopResponse{}, err
	}
	resp := &notificationservicepb.NoopResponse{NoopString: noop.NoopString}

	return resp, nil
}
