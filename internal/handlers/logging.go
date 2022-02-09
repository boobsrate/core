package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type LoggingMiddleware struct {
	logger *otelzap.Logger
}

func NewLoggingMiddleware(logger *zap.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: otelzap.New(logger.Named("server_log"))}
}

func (a *LoggingMiddleware) Apply(router *mux.Router) {
	router.Use(a.Handle)
}

func (a *LoggingMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.logger.Ctx(r.Context()).Info("access",
			zap.String("method", r.Method),
			zap.String("host", r.Host),
			zap.String("uri", r.RequestURI),
			zap.String("remote", r.RemoteAddr),
			zap.String("user-agent", r.UserAgent()),
			zap.String("referer", r.Referer()),
			zap.String("forwarded", r.Header.Get("X-Forwarded-For")),
		)
		next.ServeHTTP(w, r)
	})
}
