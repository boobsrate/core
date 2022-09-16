module github.com/boobsrate/core

go 1.18

require (
	github.com/boobsrate/apis v0.0.1
	github.com/caarlos0/env/v6 v6.10.1
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/gojuno/minimock/v3 v3.0.10
	github.com/golang-migrate/migrate/v4 v4.15.2
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/hashicorp/golang-lru v0.5.1 // indirect
	github.com/jung-kurt/gofpdf v1.0.3-0.20190309125859-24315acbbda5 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/lib/pq v1.10.7
	github.com/minio/minio-go/v7 v7.0.36
	github.com/oklog/ulid/v2 v2.1.0
	github.com/prometheus/client_golang v1.13.0
	github.com/rs/cors v1.8.2
	github.com/slok/go-http-metrics v0.10.0
	github.com/uptrace/bun v1.1.8
	github.com/uptrace/bun/dialect/pgdialect v1.1.8
	github.com/uptrace/bun/driver/pgdriver v1.1.8
	github.com/uptrace/bun/extra/bunotel v1.1.8
	github.com/uptrace/opentelemetry-go-extra/otelzap v0.1.16
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.35.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.35.0
	go.opentelemetry.io/contrib/propagators/jaeger v1.10.0
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/exporters/jaeger v1.10.0
	go.opentelemetry.io/otel/sdk v1.10.0
	go.uber.org/zap v1.23.0
	golang.org/x/exp v0.0.0-20200224162631-6cc2880d07d6 // indirect
	google.golang.org/grpc v1.49.0
)
