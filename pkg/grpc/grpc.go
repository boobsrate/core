package grpc

import (
	"context"
	"fmt"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	defaultServerKeepaliveTime     = time.Minute * 5
	defaultServerKeepaliveTimeout  = time.Second * 20
	defaultServerKeepaliveAge      = time.Minute * 5
	defaultServerKeepaliveAgeGrace = time.Minute * 1

	defaultClientDialTimeout = time.Second * 10
)

// Server can register itself to the gRPC server.
type Server interface {
	Register(srv *grpc.Server)
}

func NewGrpcServer(services []Server) *grpc.Server {
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpc_prometheus.UnaryServerInterceptor,
		otelgrpc.UnaryServerInterceptor(),
	}

	srv := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(unaryInterceptors...),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge:      defaultServerKeepaliveAge,
			MaxConnectionAgeGrace: defaultServerKeepaliveAgeGrace,
			Time:                  defaultServerKeepaliveTime,
			Timeout:               defaultServerKeepaliveTimeout,
		}),
	)

	for _, reg := range services {
		reg.Register(srv)
	}

	grpc_prometheus.Register(srv)
	grpc_prometheus.EnableHandlingTimeHistogram()

	return srv
}

func NewGrpcClient(addr string) (*grpc.ClientConn, error) {
	unaryInterceptors := []grpc.UnaryClientInterceptor{
		grpc_prometheus.UnaryClientInterceptor,
		otelgrpc.UnaryClientInterceptor(),
	}

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		grpc.WithChainUnaryInterceptor(unaryInterceptors...),
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultClientDialTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("create grpc client connection: %v", err)
	}
	return conn, err
}
