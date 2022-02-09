package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Configuration struct {
	Server   ServerConfiguration
	Database DatabaseConfig
	Tracing  TracingConfig
	Metrics  MetricsConfig
	Minio    MinioConfig
}

type ServerConfiguration struct {
	GRPCPort int `env:"GRPC_PORT" envDefault:"8081"`
	HTTPPort int `env:"HTTP_PORT" envDefault:"8088"`
}

type MinioConfig struct {
	Endpoint  string `env:"MINIO_ENDPOINT" envDefault:"storage.ops.boobsrate.com"`
	AccessKey string `env:"MINIO_ACCESS_KEY" envDefault:"golangbackend"`
	SecretKey string `env:"MINIO_SECRET_KEY" envDefault:"142701foobar"`
	Bucket    string `env:"MINIO_BUCKET" required:"true" envDefault:"tits"`
	UseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"true"`
}

type MetricsConfig struct {
	Port int `env:"METRICS_PORT" envDefault:"9000"`
}

type DatabaseConfig struct {
	Host     string `env:"DATABASE_HOST" required:"true" envDefault:"localhost"`
	Port     int    `env:"DATABASE_PORT" required:"true" envDefault:"5432"`
	User     string `env:"DATABASE_USER" required:"true" envDefault:"tits"`
	Password string `env:"DATABASE_PASSWORD" required:"true" envDefault:"tits"`
	Name     string `env:"DATABASE_NAME" required:"true" envDefault:"tits"`
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.User, d.Password, d.Host, d.Port, d.Name)
}

type TracingConfig struct {
	ProviderEndpoint string `env:"TRACING_ENDPOINT" required:"true" envDefault:"https://tempo.jaeger.ops.boobsrate.com/api/traces"`
	TracerName       string `env:"TRACING_TRACER_NAME" required:"true" envDefault:"tits"`
}

func LoadConfiguration() (*Configuration, error) {
	var config Configuration
	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
