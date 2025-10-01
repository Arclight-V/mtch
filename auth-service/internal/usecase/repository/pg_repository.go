package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Arclight-V/mtch/auth-service/internal/domain"
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
)

type TokenRepository struct {
	db *sqlx.DB
}

// NewTokenRepository returns a new TokenRepository
func NewTokenRepository(db *sqlx.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (t TokenRepository) InsertIssue(ctx context.Context, v domain.VerifyTokenIssue) error {
	//TODO implement me
	tx, err := t.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	if _, err := tx.ExecContext(ctx, insertIssueSql, v.JTI, v.UserID, v.ExpiresAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errors.New("verify token jti already exists")
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil

}
