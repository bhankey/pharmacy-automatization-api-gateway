package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bhankey/go-utils/pkg/grpc/interceptors"
	"github.com/bhankey/go-utils/pkg/logger"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/app/container"
	configinternal "github.com/bhankey/pharmacy-automatization-api-gateway/internal/config"
	"github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/middleware"
	v1 "github.com/bhankey/pharmacy-automatization-api-gateway/internal/delivery/http/v1"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	server    *http.Server
	container *container.Container
	logger    logger.Logger
}

const shutDownTimeoutSeconds = 10

// nolint: funlen, nolintlint
func NewApp(configPath string) (*App, error) {
	config, err := configinternal.GetConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init app because of config error: %w", err)
	}

	log, err := logger.GetLogger(config.Logger.Level)
	if err != nil {
		return nil, fmt.Errorf("failed to init logger error: %w", err)
	}

	log.Info("try to init data source resource")
	dataSources, err := newDataSource(config) // TODO remove dataSource struct
	if err != nil {
		return nil, err
	}

	errorHandlintInterceptor := interceptors.NewErrorHandlingInterceptor(log)

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(errorHandlintInterceptor.ClientInterceptor()),
	}

	userServiceConn, err := grpc.Dial(config.Services.User.Addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to user service: %w: ", err)
	}

	pharmacyServiceConn, err := grpc.Dial(config.Services.Pharmacy.Addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to establish connection to pharmacy service: %w: ", err)
	}

	dependencies := container.NewContainer(
		log,
		dataSources.db,
		dataSources.db,
		dataSources.redisClient,
		config.Secure.JwtKey,
		userServiceConn, pharmacyServiceConn,
	)

	mainRouter := chi.NewRouter()

	mainRouter.Use(func(handler http.Handler) http.Handler {
		return middleware.LoggingMiddleware(log)(handler)
	})
	mainRouter.Use(middleware.FingerPrint)

	v1Router := v1.NewRouter(
		dependencies.GetV1SwaggerHandler(),
		dependencies.GetV1AuthHandler(),
		dependencies.GetV1UserHandler(),
		dependencies.GetV1PharmacyHandler(),
	)

	mainRouter.Mount("/v1", v1Router)

	server := &http.Server{
		Addr:    ":" + config.Server.Port,
		Handler: mainRouter,
	}

	return &App{logger: log, server: server, container: dependencies}, nil
}

func (a *App) Start() {
	a.logger.Info("staring server on port: " + a.server.Addr)
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			a.logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	a.logger.Info("received signal to shutdown server")
	ctx, cancel := context.WithTimeout(context.Background(), shutDownTimeoutSeconds*time.Second)
	defer cancel()
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error(err)
	}

	<-ctx.Done()

	a.container.CloseAllConnections()

	a.logger.Info("server was shutdown")
}
