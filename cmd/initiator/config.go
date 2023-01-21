package main

import (
	"github.com/caarlos0/env/v6"
)

// Configuration represents application configuration for serve action.
type Configuration struct {
	Database DatabaseConfig
	Minio    MinioConfig
	Images   ImagesConfig
}

type ImagesConfig struct {
	OptimizerEndpoint string `env:"IMAGES_OPTIMIZER_ENDPOINT" envDefault:"https://img.optimizer.akuzyashin.pw"`
}

type MinioConfig struct {
	Endpoint  string `env:"MINIO_ENDPOINT" envDefault:"s3.rate-tits.online:443"`
	AccessKey string `env:"MINIO_ACCESS_KEY" envDefault:""`
	SecretKey string `env:"MINIO_SECRET_KEY" envDefault:""`
	Bucket    string `env:"MINIO_BUCKET" envDefault:"boobsrate"`
	UseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"true"`
}

type DatabaseConfig struct {
	DatabaseDSN string `env:"DATABASE_DSN"`
}

// LoadConfiguration returns a new application configuration parsed from environment variables.
func LoadConfiguration() (*Configuration, error) {
	var config Configuration
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
