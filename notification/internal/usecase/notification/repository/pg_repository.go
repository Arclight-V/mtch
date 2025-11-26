package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"

	domain "github.com/Arclight-V/mtch/notification/internal/domain/notification"
)

// NotificationRepository is a wrapper around *sqlx.DB
type NotificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) InsertIssue(ctx context.Context, vc *domain.VerificationCode) error {
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
