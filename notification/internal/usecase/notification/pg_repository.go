package notification

import (
	"context"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
)

// PGRepository interface
//
//go:generate mockgen -source=$GOFILE -package=mocks -destination=./mocks/pg_repository_mock.go
type PGRepository interface {
	InsertIssue(ctx context.Context, vc *domain.VerificationCode) error
}
