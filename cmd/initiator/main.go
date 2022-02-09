package main

import (
	"fmt"
	"log"

	"github.com/boobsrate/core/internal/applications/initiator"
	"github.com/boobsrate/core/internal/repository/postgres"
	"github.com/boobsrate/core/internal/services/tits"
	storage "github.com/boobsrate/core/internal/storage/minio"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Errorf("creating logger: %w", err))
	}
	logger = logger.Named("initiator")
	defer logger.Sync() // nolint: errcheck

	cfg, err := LoadConfiguration()
	if err != nil {
		log.Fatal("load configuration: ", zap.Error(err))
	}

	pgDB := postgres.NewPostgresDatabase(cfg.Database.DSN())
	titsRepo := postgres.NewTitsRepository(pgDB)

	minioClient, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKey, cfg.Minio.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		logger.Fatal("creating minio client: ", zap.Error(err))
	}
	titsStorage := storage.NewMinioStorage(minioClient, cfg.Minio.Bucket)
	titsService := tits.NewService(titsRepo, titsStorage, logger, nil)

	initiatorApp := initiator.NewService(logger, titsService)
	initiatorApp.Run()
}
