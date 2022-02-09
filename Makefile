.EXPORT_ALL_VARIABLES:
GOBIN = $(shell pwd)/bin
GOFLAGS = -mod=vendor
GO111MODULE = on
SHELL=/bin/bash

.PHONY: deps
deps:
	@go mod download
	@go mod vendor
	@go mod tidy

.PHONY: tools
tools: deps
	@go install github.com/gojuno/minimock/v3/cmd/minimock
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate

.PHONY: generate
generate: tools
	@export PATH=$(shell pwd)/bin:$(PATH); go generate ./...


.PHONY: test
test:
	@go test ./...


.PHONY: lint
lint:
	@golangci-lint run

.PHONY: clean
clean:
	@rm -fv ./bin/*

.PHONY: migration
migration:
	@./bin/migrate create -ext sql -dir migrations -seq -digits 8 $(NAME)
