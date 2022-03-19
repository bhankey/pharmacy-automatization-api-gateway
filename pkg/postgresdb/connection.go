package postgresdb

import (
	"fmt"

	// driver for connection.
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

// NewClient connect to postgres and ping it. Return connection to database.
func NewClient(host, port, user, password, dbName string) (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	db, err := sqlx.Connect("pgx", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres client error: %w", err)
	}

	return db, nil
}
