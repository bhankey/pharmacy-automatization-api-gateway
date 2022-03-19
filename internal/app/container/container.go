package container

import (
	"github.com/bhankey/pharmacy-automatization-api-gateway/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type Container struct {
	masterPostgresDB *sqlx.DB
	slavePostgresDB  *sqlx.DB
	redisConnection  *redis.Client
	logger           logger.Logger

	userServiceConn     *grpc.ClientConn
	pharmacyServiceConn *grpc.ClientConn

	jwtKey          string
	smtpMessageFrom string

	dependencies map[string]interface{}
}

func NewContainer(
	log logger.Logger,
	masterPostgres, slavePostgres *sqlx.DB,
	redis *redis.Client,
	jwtKey string,
	userServiceConn, pharmacyServiceConn *grpc.ClientConn,
) *Container {
	return &Container{
		masterPostgresDB:    masterPostgres,
		slavePostgresDB:     slavePostgres,
		redisConnection:     redis,
		logger:              log,
		userServiceConn:     userServiceConn,
		pharmacyServiceConn: pharmacyServiceConn,
		jwtKey:              jwtKey,
		dependencies:        make(map[string]interface{}),
	}
}

func (c *Container) CloseAllConnections() {
	if err := c.masterPostgresDB.Close(); err != nil {
		c.logger.Errorf("failed to close master postgres connection error: %v", err)
	}

	if err := c.slavePostgresDB.Close(); err != nil {
		c.logger.Errorf("failed to close slave postgres connection error: %v", err)
	}

	if err := c.redisConnection.Close(); err != nil {
		c.logger.Errorf("failed to close redis connection error: %v", err)
	}
}
