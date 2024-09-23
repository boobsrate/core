FROM golang:1.20.10-alpine3.17 AS builder
RUN apk update && apk add --no-cache ca-certificates git gcc make libc-dev binutils-gold
WORKDIR tits
COPY . .
RUN go mod tidy && go get -d -v

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -trimpath -o /bin/migrator ./cmd/migrate
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -trimpath -o /bin/parser ./cmd/parser
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -trimpath -o /bin/tits ./cmd/tits


FROM scratch
COPY --from=builder /bin/migrator /bin/migrator
COPY --from=builder /bin/tits /bin/tits
COPY --from=builder /bin/parser /bin/parser
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY migrations migrations
COPY assets assets

ENV MIGRATIONS_DIR /migrations
