package testingpostgres

import (
	"context"
	"errors"
	"log"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/golang-migrate/migrate/v4"
	migrage_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func ApplyMigrations(t *testing.T, dsn, sourceURL string) {
	t.Helper()

	db, err := sqlx.Open("pgx", dsn)
	require.NoError(t, err)
	defer db.Close()

	driver, err := migrage_postgres.WithInstance(db.DB, &migrage_postgres.Config{})
	require.NoError(t, err)

	m, err := migrate.NewWithDatabaseInstance(
		sourceURL, "postgres", driver)
	require.NoError(t, err)

	err = m.Up()
	if err != nil || errors.Is(err, migrate.ErrNoChange) {
		t.Fatalf("migrate up %v", err)
	}
}

// CreatingPostgresContainer Run creates an instance of the Postgres container type from image
func CreatingPostgresContainer(t *testing.T, img string) (*postgres.PostgresContainer, error, func()) {

	ctx := context.Background()

	dbName := "mtch_db_test"
	dbUser := "mtch_db_test_user"
	dbPassword := "password"

	postgresContainer, err := postgres.Run(ctx,
		img,
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return nil, err, nil
	}

	testcontainers.CleanupContainer(t, postgresContainer)
	require.NoError(t, err)

	terminateContainerFunc := func() {
		if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
			log.Printf("failed to terminate container: %s", err)
		}
	}

	return postgresContainer, nil, terminateContainerFunc
}
