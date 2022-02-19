package config

import (
	"github.com/caarlos0/env/v6"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfig
	Tracing  TracingConfig
	Metrics  MetricsConfig
	Minio    MinioConfig
	Images   ImagesConfig
}

type ServerConfiguration struct {
	GRPCPort int `env:"GRPC_PORT" envDefault:"8081"`
	HTTPPort int `env:"HTTP_PORT" envDefault:"8088"`
}

type ImagesConfig struct {
	PublicEndpoint    string `env:"IMAGES_PUBLIC_ENDPOINT" envDefault:"https://s3.boobsrate.com"`
	OptimizerEndpoint string `env:"IMAGES_OPTIMIZER_ENDPOINT" envDefault:"http://image-optimizer.images:3000"`
}

type MinioConfig struct {
	Endpoint  string `env:"MINIO_ENDPOINT" envDefault:"minio.images:9000"`
	AccessKey string `env:"MINIO_ACCESS_KEY" envDefault:""`
	SecretKey string `env:"MINIO_SECRET_KEY" envDefault:""`
	Bucket    string `env:"MINIO_BUCKET" envDefault:"tits"`
	UseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"true"`
}

type MetricsConfig struct {
	Port int `env:"METRICS_PORT" envDefault:"9090"`
}

type DatabaseConfig struct {
	DatabaseDSN string `env:"DATABASE_DSN" required:"true"`
}

type TracingConfig struct {
	ProviderEndpoint string `env:"TRACING_ENDPOINT" required:"true" envDefault:"http://monitoring-tempo.monitoring:14268/api/traces"`
	TracerName       string `env:"TRACING_TRACER_NAME" required:"true" envDefault:"tits"`
}

func LoadConfiguration() (*Configuration, error) {
	var config Configuration
	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
