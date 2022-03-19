package app

import (
	"context"
	"fmt"

	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/config"
	"github.com/bhankey/pharmacy-automatization-api-gateway/pkg/postgresdb"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type dataSources struct {
	db          *sqlx.DB
	redisClient *redis.Client
}

func newDataSource(config config.Config) (*dataSources, error) {
	postgresDB, err := postgresdb.NewClient(
		config.Postgres.Host,
		config.Postgres.Port,
		config.Postgres.User,
		config.Postgres.Password,
		config.Postgres.DBName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to init postgres connection error: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       0,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("failed init redis conntection: %w", err)
	}

	return &dataSources{
		db:          postgresDB,
		redisClient: rdb,
	}, nil
}
