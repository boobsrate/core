FROM golang:1.18-alpine3.15 AS builder
RUN apk update && apk add --no-cache ca-certificates git gcc make libc-dev binutils-gold
WORKDIR tits
COPY . .
RUN go mod tidy && go get -d -v

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linuxgo build -o /bin/migrator ./cmd/migrate
#RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -trimpath -o /bin/initiator ./cmd/initiator
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o /bin/tits ./cmd/tits


FROM scratch
COPY --from=builder /bin/migrator /bin/migrator
COPY --from=builder /bin/tits /bin/tits
#COPY --from=builder /bin/initiator /bin/initiator
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY migrations migrations
##COPY assets assets

ENV MIGRATIONS_DIR /migrations
