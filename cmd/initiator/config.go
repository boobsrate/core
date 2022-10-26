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
	OptimizerEndpoint string `env:"IMAGES_OPTIMIZER_ENDPOINT" envDefault:"http://image-optimizer.images:3000"`
}

type MinioConfig struct {
	Endpoint  string `env:"MINIO_ENDPOINT" envDefault:"minio.images:9000"`
	AccessKey string `env:"access_key" envDefault:""`
	SecretKey string `env:"secret_key" envDefault:""`
	Bucket    string `env:"MINIO_BUCKET" envDefault:"tits"`
	UseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"false"`
}

type DatabaseConfig struct {
	DatabaseDSN string `env:"connection_dsn"`
}

// LoadConfiguration returns a new application configuration parsed from environment variables.
func LoadConfiguration() (*Configuration, error) {
	var config Configuration
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
