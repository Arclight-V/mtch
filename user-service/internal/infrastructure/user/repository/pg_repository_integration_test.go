package repository

import (
	"context"
	"errors"
	domain "github.com/Arclight-V/mtch/user-service/internal/domain/user"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migrage_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func applyMigrations(t *testing.T, dsn string) {
	t.Helper()

	db, err := sqlx.Open("pgx", dsn)
	require.NoError(t, err)
	defer db.Close()

	driver, err := migrage_postgres.WithInstance(db.DB, &migrage_postgres.Config{})
	require.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+filepath.Join("..", "..", "..", "..", "..", "db", "migrations"), "postgres", driver)
	require.NoError(t, err)

	err = m.Up()
	if err != nil || errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("migrate up %v", err)
	}

}

func creatingPostgresContainer(t *testing.T) (*postgres.PostgresContainer, error) {

	ctx := context.Background()

	dbName := "mtch_db_test"
	dbUser := "mtch_db_test_user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:18",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, err
	}

	testcontainers.CleanupContainer(t, postgresContainer)
	require.NoError(t, err)
	return postgresContainer, nil
}

func TestCreate_Integration(t *testing.T) {
	postgresContainer, err := creatingPostgresContainer(t)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}()

	ctx := context.Background()
	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	applyMigrations(t, dsn)

	//err = postgresContainer.Snapshot(ctx)
	//require.NoError(t, err)

	personalData := domain.PersonalData{
		FirstName:    "John",
		LastName:     "Doe",
		Contact:      "email",
		Phone:        "+7999999999",
		Email:        "a@b.com",
		Password:     "password",
		DateBirthday: time.Date(1992, time.Month(11), 28, 0, 0, 0, 0, time.UTC),
		Gender:       domain.Male,
	}

	pendingUser, err := domain.NewPendingUser(&personalData)

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
