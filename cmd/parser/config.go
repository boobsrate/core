package main

import (
	"github.com/caarlos0/env/v6"
)

// Configuration represents application configuration for serve action.
type Configuration struct {
	Base     BaseConfig
	Database DatabaseConfig
	Minio    MinioConfig
	Images   ImagesConfig
	Proxy    ProxyConfig
}

type BaseConfig struct {
	WithFill bool `env:"WITH_FILL" envDefault:"true"`
}

type ProxyConfig struct {
	ProxyEndpointEU  string `env:"PROXY_ENDPOINT_EU" envDefault:"http://opera-proxy-eu.proxypool:8080"`
	ProxyEndpointAll string `env:"PROXY_ENDPOINT_ALL" envDefault:"http://opera-proxy-all.proxypool:8080"`
}

type ImagesConfig struct {
	OptimizerEndpoint string `env:"IMAGES_OPTIMIZER_ENDPOINT" envDefault:"https://img.optimizer.akuzyashin.pw"`
}

type MinioConfig struct {
	Endpoint  string `env:"MINIO_ENDPOINT" envDefault:"minio.minio:9000"`
	AccessKey string `env:"MINIO_ACCESS_KEY" envDefault:""`
	SecretKey string `env:"MINIO_SECRET_KEY" envDefault:""`
	Bucket    string `env:"MINIO_BUCKET" envDefault:"tits"`
	UseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"false"`
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
