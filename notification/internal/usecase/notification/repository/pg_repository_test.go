package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"github.com/Arclight-V/mtch/notification/internal/infrastructure/codegen"
)

func TestNotificationRepositoryInsertIsue(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	defer sqlxDB.Close()

	pgRepository := NewNotificationRepository(sqlxDB)
	codegen := codegen.NewNoopCodeGenerator()
	ui := uuid.New()
	mockVerificationCode := codegen.NewVerificationCode(ui.String())

	columns := []string{"user_id", "code_hash", "expires_at", "purpose", "attempts", "max_attempts"}
	_ = sqlmock.NewRows(columns).AddRow(
		ui,
		mockVerificationCode.Code,
		mockVerificationCode.ExpiresAt,
		mockVerificationCode.Purpose,
		mockVerificationCode.Attempts,
		mockVerificationCode.MaxAttempts,
	)

	mock.ExpectExec(insertIssueSql).WithArgs(
		ui,
		mockVerificationCode.Code,
		mockVerificationCode.ExpiresAt,
		mockVerificationCode.Purpose,
		mockVerificationCode.Attempts,
		mockVerificationCode.MaxAttempts,
	).WillReturnResult(sqlmock.NewResult(1, 1))

	err = pgRepository.InsertIssue(context.Background(), mockVerificationCode)
	require.NoError(t, err)
}
