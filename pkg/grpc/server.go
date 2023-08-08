package grpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GracefulServer struct {
	addr   string
	server *grpc.Server
	log    *zap.Logger

	wg   sync.WaitGroup
	dead chan struct{}
}

func NewGracefulServer(port int, server *grpc.Server, log *zap.Logger) *GracefulServer {
	return &GracefulServer{
		addr:   fmt.Sprintf(":%d", port),
		log:    log.Named("grpc_server"),
		server: server,
		dead:   make(chan struct{}),
	}
}

func (s *GracefulServer) Serve() error {
	s.log.Info("Server starting...")
	defer s.log.Info("Server started")
	s.wg.Add(1)
	go s.startServer()
	return nil
}

func (s *GracefulServer) startServer() {
	defer s.wg.Done()
	defer close(s.dead)

	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		s.handleListenError(err)
		return
	}
	defer lis.Close() // nolint

	if err := s.server.Serve(lis); err != nil {
		s.handleServeError(err)
	}
}

func (s *GracefulServer) handleListenError(err error) {
	s.log.Error(fmt.Sprintf("listen: %v", err))
}

func (s *GracefulServer) handleServeError(err error) {
	s.log.Error(fmt.Sprintf("serve: %v", err))
}

func (s *GracefulServer) Shutdown(ctx context.Context) error {
	s.log.Info("Server stopping...")
	defer s.log.Info("Server stopped")
	shutdown := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(shutdown)
	}()
	select {
	case <-ctx.Done():
		s.server.Stop()
		<-shutdown
	case <-shutdown:
	}
	s.wg.Wait()
	return nil
}

func (s *GracefulServer) Dead() <-chan struct{} {
	return s.dead
}