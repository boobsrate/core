module github.com/boobsrate/core

go 1.16

require (
	github.com/boobsrate/apis v0.0.1
	github.com/caarlos0/env/v6 v6.9.1
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/godbus/dbus v0.0.0-20190422162347-ade71ed3457e // indirect
	github.com/gojuno/minimock/v3 v3.0.10
	github.com/golang-migrate/migrate/v4 v4.15.1
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/jung-kurt/gofpdf v1.0.3-0.20190309125859-24315acbbda5 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/lib/pq v1.10.0
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/minio/minio-go/v7 v7.0.21
	github.com/oklog/ulid/v2 v2.0.2
	github.com/prometheus/client_golang v1.12.1
	github.com/rs/cors v1.8.2
	github.com/slok/go-http-metrics v0.10.0
	github.com/uptrace/bun v1.0.22
	github.com/uptrace/bun/dialect/pgdialect v1.0.22
	github.com/uptrace/bun/driver/pgdriver v1.0.22
	github.com/uptrace/bun/extra/bunotel v1.0.22
	github.com/uptrace/opentelemetry-go-extra/otelzap v0.1.8
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.29.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.28.0
	go.opentelemetry.io/contrib/propagators/jaeger v1.4.0
	go.opentelemetry.io/otel v1.4.0
	go.opentelemetry.io/otel/exporters/jaeger v1.3.0
	go.opentelemetry.io/otel/sdk v1.3.0
	go.uber.org/zap v1.20.0
	golang.org/x/exp v0.0.0-20200224162631-6cc2880d07d6 // indirect
	google.golang.org/grpc v1.44.0
)
