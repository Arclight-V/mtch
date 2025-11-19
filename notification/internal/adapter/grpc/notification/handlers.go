package notification

import (
	"context"
	"log"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
)

func (s *notificationServiceServer) NoopMethod(ctx context.Context, req *notificationservicepb.NoopRequest) (*notificationservicepb.NoopResponse, error) {
	log.Println("NoopMethod called:", req.NoopString)

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
