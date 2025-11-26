package postgres

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"

	"github.com/Arclight-V/mtch/pkg/platform/config"
)

// NewPsqlDB returns instance *sqlx.DB
func NewPsqlDB(c *config.PostgresCfg, opts ...Option) (*sqlx.DB, error) {

	options := options{}
	for _, o := range opts {
		o.apply(&options)
	}

	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		c.Host,
		c.Port,
		c.User,
		c.DBName,
		c.Password,
	)

	var (
		db  *sqlx.DB
		err error
	)

	if options.context != nil {
		db, err = sqlx.ConnectContext(options.context, c.PgDriver, dataSourceName)

	} else {
		db, err = sqlx.Connect(c.PgDriver, dataSourceName)
	}
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(options.maxOpenConns)
	db.SetConnMaxLifetime(options.connMaxLifetime)
	db.SetMaxIdleConns(options.maxIdleConns)
	db.SetConnMaxIdleTime(options.connMaxIdleTime)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
