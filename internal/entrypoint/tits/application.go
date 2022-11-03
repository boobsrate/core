package tits

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/boobsrate/core/internal/applications/abyss"
	"github.com/boobsrate/core/internal/applications/websockethub"
	"github.com/boobsrate/core/internal/config"
	titsproto "github.com/boobsrate/core/internal/grpcapi/tits"
	authhandlers "github.com/boobsrate/core/internal/handlers/auth"
	titshandlers "github.com/boobsrate/core/internal/handlers/tits"
	wshandler "github.com/boobsrate/core/internal/handlers/websocket"
	"github.com/boobsrate/core/internal/repository/postgres"
	titssvc "github.com/boobsrate/core/internal/services/tits"
	minio2 "github.com/boobsrate/core/internal/storage/minio"
	"github.com/boobsrate/core/pkg/grpc"
	"github.com/boobsrate/core/pkg/logging"
	"github.com/boobsrate/core/pkg/migrations"
	"github.com/boobsrate/core/pkg/observer"
	"github.com/boobsrate/core/pkg/server"
	"github.com/boobsrate/core/pkg/tracing"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	otelmiddleware "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.uber.org/zap"
)

func Run() error {
	logger, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("failed to create logger: %v", err)
	}
	defer logger.Sync() // nolint:errcheck
	appLogger := logging.NewLogger(logger, "tits")

	cfg, err := config.LoadConfiguration()
	if err != nil {
		appLogger.Error("loading configuration", zap.Error(err))
		return fmt.Errorf("loading configuration: %v", err)
	}

	migrationManager, err := migrations.NewManager(cfg.Database.DatabaseDSN)
	if err != nil {
		appLogger.Error("create migrations manager", zap.Error(err))
	}

	migrateCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err = migrationManager.Wait(migrateCtx)
	if err != nil {
		appLogger.Error("wait migrations", zap.Error(err))
	}

	logging.SetInternalGRPCLogger(logger.Named("grpc_logger"))

	healthMetricsHandler := tracing.NewGracefulMetricsServer()
	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Metrics.Port),
		Handler: healthMetricsHandler,
	}

	httpMetricsServer := server.NewGracefulServer(metricsServer, logger.Named("metrics_server"))

	rootRouter := mux.NewRouter()

	tp, err := tracing.NewTracingProvider(cfg.Tracing.ProviderEndpoint, cfg.Tracing.TracerName)
	if err != nil {
		appLogger.Error("create tracing provider", zap.Error(err))
		return fmt.Errorf("creating tracing provider: %v", err)
	}

	rootRouter.Use(otelmiddleware.Middleware("tits"))

	wsHub := websockethub.NewWebsocketsHub(logger)
	wsHandler := wshandler.NewWebsocketHandler(logger, wsHub.ClientsChannel())
	wsHandler.Register(rootRouter)

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: cfg.Minio.UseSSL,
	})
	if err != nil {
		appLogger.Error("creating minio client", zap.Error(err))
		return fmt.Errorf("creating minio client: %v", err)
	}

	minioStorage := minio2.NewMinioStorage(minioClient, cfg.Minio.Bucket, cfg.Images.PublicEndpoint)

	database := postgres.NewPostgresDatabase(cfg.Database.DatabaseDSN)
	titsRepo := postgres.NewTitsRepository(database)

	titsService := titssvc.NewService(titsRepo, minioStorage, logger, wsHub.MessagesChannel(), cfg.Images.OptimizerEndpoint)
	titsGrpcServer := titsproto.NewTitsGRPCServer(titsService)

	titsHttpService := titshandlers.NewTitsHandler(titsService)
	titsHttpService.Register(rootRouter)

	authhandler := authhandlers.NewAuthHandler()
	authhandler.Register(rootRouter)

	grpcServer := grpc.NewGrpcServer([]grpc.Server{
		titsGrpcServer,
	})

	gracefulServer := grpc.NewGracefulServer(cfg.Server.GRPCPort, grpcServer, logger)

	rootServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.HTTPPort),
		Handler: tracing.ApplyPrometheusMiddleware(server.ApplyCors(rootRouter), "titsbackend"),
	}

	httpRootServer := server.NewGracefulServer(rootServer, logger.Named("http_server"))

	abyssKeeper := abyss.NewKeeper(logger, titsService)

	obs := observer.NewObserver()

	obs.AddOpener(observer.OpenerFunc(func() error {
		return gracefulServer.Serve()
	}))

	obs.AddOpener(observer.OpenerFunc(func() error {
		return httpMetricsServer.Serve()
	}))

	obs.AddOpener(observer.OpenerFunc(func() error {
		return httpRootServer.Serve()
	}))

	obs.AddContextCloser(observer.ContextCloserFunc(func(ctx context.Context) error {
		return gracefulServer.Shutdown(ctx)
	}))

	obs.AddContextCloser(observer.ContextCloserFunc(func(ctx context.Context) error {
		return tp.Shutdown(ctx)
	}))

	obs.AddContextCloser(observer.ContextCloserFunc(func(ctx context.Context) error {
		return httpRootServer.Shutdown(ctx)
	}))

	obs.AddContextCloser(observer.ContextCloserFunc(func(ctx context.Context) error {
		return httpMetricsServer.Shutdown(ctx)
	}))

	obs.AddUpper(func(ctx context.Context) {
		wsHub.Run(ctx)
	})

	obs.AddUpper(func(ctx context.Context) {
		abyssKeeper.Run(ctx)
	})

	obs.AddUpper(func(ctx context.Context) {
		select {
		case <-ctx.Done():
		case <-gracefulServer.Dead():
		case <-httpRootServer.Dead():
		case <-httpMetricsServer.Dead():
		case <-wsHub.Dead():
		case <-abyssKeeper.Dead():
		}
	})

	return obs.Run()

}
