//go:build integration

package repository

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	tesging_postgres "github.com/Arclight-V/mtch/pkg/testing/postgres"

	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
	user_test_data "github.com/Arclight-V/mtch/user-service/internal/domain/user/testdata"
)

func TestCreate_Integration(t *testing.T) {
	img := "postgres:18"
	postgresContainer, err, terminateContainerFunc := tesging_postgres.CreatingPostgresContainer(t, img)
	if err != nil {
		t.Fatal(err)
	}
	defer terminateContainerFunc()

	ctx := context.Background()
	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	tesging_postgres.ApplyMigrations(t, dsn, "file://"+filepath.Join("..", "..", "..", "..", "..", "db", "migrations"))

	//err = postgresContainer.Snapshot(ctx)
	//require.NoError(t, err)

	pendingUser, err := user_test_data.NewTestPendingUser()

	t.Run("Test creating a user", func(t *testing.T) {
		t.Cleanup(func() {
			//err = postgresContainer.Restore(ctx)
			require.NoError(t, err)
		})
		conn, err := pgx.Connect(ctx, dsn)
		require.NoError(t, err)
		defer conn.Close(ctx)

		if err != nil {
			t.Fatal("failed to create user", err)
		}

		var user domain.User
		if err := conn.QueryRow(ctx, createPendingUserQuery,
			pendingUser.UserID,
			pendingUser.FirstName,
			pendingUser.LastName,
			pendingUser.Contact,
			pendingUser.Phone,
			pendingUser.Email,
			pendingUser.Password,
			pendingUser.DateBirthday,
			pendingUser.Gender,
			pendingUser.Role,
		).Scan(
			&user.UserID,
			&user.FirstName,
			&user.LastName,
			&user.Contact,
			&user.Phone,
			&user.Email,
			&user.Password,
			&user.DateBirthday,
			&user.Gender,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			t.Fatal("failed to create user", err)
		}

		require.NoError(t, err)

		require.Equal(t, pendingUser.FirstName, user.FirstName)
		require.Equal(t, pendingUser.LastName, user.LastName)
		require.Equal(t, pendingUser.Contact, user.Contact)
		require.Equal(t, pendingUser.Phone, user.Phone)
		require.Equal(t, pendingUser.Email, user.Email)
		require.Equal(t, pendingUser.Password, user.Password)
		require.Equal(t, pendingUser.DateBirthday, user.DateBirthday)
		require.Equal(t, pendingUser.Role, user.Role)
		require.Equal(t, pendingUser.Activated, user.Activated)

	})

	t.Run("Test not creating a user", func(t *testing.T) {
		t.Cleanup(func() {
			//err = postgresContainer.Restore(ctx)
			require.NoError(t, err)
		})
		conn, err := pgx.Connect(ctx, dsn)
		require.NoError(t, err)
		defer conn.Close(ctx)

		if err != nil {
			t.Fatal("failed to create user", err)
		}

		var user domain.User
		err = conn.QueryRow(ctx, createPendingUserQuery,
			pendingUser.UserID,
			pendingUser.FirstName,
			pendingUser.LastName,
			pendingUser.Contact,
			pendingUser.Phone,
			pendingUser.Email,
			pendingUser.Password,
			pendingUser.DateBirthday,
			pendingUser.Gender,
			pendingUser.Role,
		).Scan(
			&user.UserID,
			&user.FirstName,
			&user.LastName,
			&user.Contact,
			&user.Phone,
			&user.Email,
			&user.Password,
			&user.DateBirthday,
			&user.Gender,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		require.Error(t, err)

	})
}
