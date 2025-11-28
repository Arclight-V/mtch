package repository

import (
	"context"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/jmoiron/sqlx"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
)

// NotificationRepoDB is a wrapper around *sqlx.DB
type NotificationRepoDB struct {
	db *sqlx.DB

	logger log.Logger
}

// NewNotificationRepoDB returns new NotificationRepoDB
func NewNotificationRepoDB(logger log.Logger, db *sqlx.DB) *NotificationRepoDB {
	return &NotificationRepoDB{
		logger: logger,
		db:     db,
	}
}

// InsertIssue inserts code to db
func (r *NotificationRepoDB) InsertIssue(ctx context.Context, vc *domain.VerificationCode) error {
	level.Debug(r.logger).Log("method", "InsertIssue", "vc", vc)

	if _, err := r.db.ExecContext(ctx, insertIssueSql,
		vc.UserID,
		vc.Code,
		vc.ExpiresAt,
		vc.Purpose,
		vc.Attempts,
		vc.MaxAttempts,
	); err != nil {
		return fmt.Errorf("insert verification code: %w", err)
	}

	return nil
}
