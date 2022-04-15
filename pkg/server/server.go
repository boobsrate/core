package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type GracefulServer struct {
	server *http.Server
	log    *otelzap.Logger

	wg   sync.WaitGroup
	dead chan struct{}
}

func NewGracefulServer(server *http.Server, log *zap.Logger) *GracefulServer {
	return &GracefulServer{
		log:    otelzap.New(log),
		server: server,
	}
}

func (s *GracefulServer) Serve() error {
	s.log.Info("Server starting...")
	defer s.log.Info("Server started")
	s.wg.Add(1)
	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			close(s.dead)
		}
		s.wg.Done()
	}()
	return nil
}

func (s *GracefulServer) Shutdown(ctx context.Context) error {
	s.log.Info("Server stopping...")
	defer s.log.Info("Server stopped")
	return s.server.Shutdown(ctx)
}

func (s *GracefulServer) Dead() <-chan struct{} {
	return s.dead
}

func ApplyCors(router *mux.Router) http.Handler {
	c := cors.New(cors.Options{
		AllowCredentials: true,
		Debug:            false,
		AllowedOrigins:   []string{"boobsrate.com", "dev.boobsrate.com"},
	})

	handler := c.Handler(router)
	return handler
}
