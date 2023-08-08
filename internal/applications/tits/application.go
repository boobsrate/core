package tits

import (
	"context"
	"fmt"
	"net/http"

	"github.com/boobsrate/core/internal/applications/websockethub"
	"github.com/boobsrate/core/internal/config"
	"github.com/boobsrate/core/internal/grpcapi/tits"
	"github.com/boobsrate/core/internal/handlers"
	tits2 "github.com/boobsrate/core/internal/handlers/tits"
	wshandler "github.com/boobsrate/core/internal/handlers/websocket"
	"github.com/boobsrate/core/internal/repository/postgres"
	"github.com/boobsrate/core/internal/services/tits"
	minio2 "github.com/boobsrate/core/internal/storage/minio"
	"github.com/boobsrate/core/pkg/grpc"
	"github.com/boobsrate/core/pkg/logging"
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

	logging.SetInternalGRPCLogger(logger.Named("grpc_logger"))

	rootRouter := mux.NewRouter()

	tp, err := tracing.NewTracingProvider(cfg.Tracing.ProviderEndpoint, cfg.Tracing.TracerName)
	if err != nil {
		appLogger.Error("create tracing provider", zap.Error(err))
		return fmt.Errorf("creating tracing provider: %v", err)
	}

	rootRouter.Use(otelmiddleware.Middleware("tits"))

	loggingMiddleware := handlers.NewLoggingMiddleware(logger)
	loggingMiddleware.Apply(rootRouter)

	wsHub, wsHandler, minioClient, minioStorage, database, titsRepo, titsService, titsGrpcServer, titsHttpService, grpcServer, gracefulServer, rootServer, httpRootServer, obs := createServices(cfg, logger, rootRouter, appLogger)

	registerHandlers(wsHub, wsHandler, titsHttpService, rootRouter)

	createServers(gracefulServer, httpRootServer, obs)

	runObserver(tp, wsHub, gracefulServer, httpRootServer, obs)

	return obs.Run()
}

func createServices(cfg *config.Config, logger *zap.Logger, rootRouter *mux.Router, appLogger *zap.Logger) (*websockethub.WebsocketsHub, *wshandler.WebsocketHandler, *minio.Client, *minio2.MinioStorage, *postgres.PostgresDatabase, *postgres.TitsRepository, *tits.Service, *titspbv1.TitsGRPCServer, *tits2.TitsHandler, *grpc.Server, *grpc.GracefulServer, *http.Server, *server.GracefulServer, *observer.Observer) {
	healthMetricsHandler := tracing.NewGracefulMetricsServer()
	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Metrics.Port),
		Handler: healthMetricsHandler,
	}

	httpMetricsServer := server.NewGracefulServer(metricsServer, logger.Named("metrics_server"))

	wsHub := websockethub.NewWebsocketsHub(logger)
	wsHandler := wshandler.NewWebsocketHandler(logger, wsHub.ClientsChannel())
	wsHandler.Register(rootRouter)

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: cfg.Minio.UseSSL,
	})
 	if err != nil {
 		appLogger.Error("creating minio client", zap.Error(err))
 		return nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, err
 	}

	minioStorage := minio2.NewMinioStorage(minioClient, cfg.Minio.Bucket, cfg.Images.PublicEndpoint)

	database := postgres.NewPostgresDatabase(cfg.Database.DatabaseDSN)
	titsRepo := postgres.NewTitsRepository(database)

	titsService := tits.NewService(titsRepo, minioStorage, logger, wsHub.MessagesChannel(), cfg.Images.OptimizerEndpoint)
	titsGrpcServer := titspbv1.NewTitsGRPCServer(titsService)

	titsHttpService := tits2.NewTitsHandler(titsService)
	titsHttpService.Register(rootRouter)

	grpcServer := grpc.NewGrpcServer([]grpc.Server{
		titsGrpcServer,
	})

	gracefulServer := grpc.NewGracefulServer(cfg.Server.GRPCPort, grpcServer, logger)

	rootServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.HTTPPort),
		Handler: tracing.ApplyPrometheusMiddleware(server.ApplyCors(rootRouter), "titsbackend"),
	}

	httpRootServer := server.NewGracefulServer(rootServer, logger.Named("http_server"))

	obs := observer.NewObserver()

	return wsHub, wsHandler, minioClient, minioStorage, database, titsRepo, titsService, titsGrpcServer, titsHttpService, grpcServer, gracefulServer, rootServer, httpRootServer, obs
}

func registerHandlers(wsHub *websockethub.WebsocketsHub, wsHandler *wshandler.WebsocketHandler, titsHttpService *tits2.TitsHandler, rootRouter *mux.Router) {
	wsHandler.Register(rootRouter)
	titsHttpService.Register(rootRouter)
}

func createServers(gracefulServer *grpc.GracefulServer, httpRootServer *server.GracefulServer, obs *observer.Observer) {
	obs.AddOpener(observer.OpenerFunc(func() error {
		return gracefulServer.Serve()
	}))

	obs.AddOpener(observer.OpenerFunc(func() error {
		return httpRootServer.Serve()
	}))
}

func runObserver(tp *tracing.TracingProvider, wsHub *websockethub.WebsocketsHub, gracefulServer *grpc.GracefulServer, httpRootServer *server.GracefulServer, obs *observer.Observer) {
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
		select {
		case <-ctx.Done():
		case <-gracefulServer.Dead():
		case <-httpRootServer.Dead():
		}
	})
}