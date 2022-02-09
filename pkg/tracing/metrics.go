package tracing

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type aliveHandler struct{}

func (h *aliveHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("i'm alive"))
	if err != nil {
		log.Printf("writing alive response: %v", err)
	}
}

func NewGracefulMetricsServer() http.Handler {
	r := http.NewServeMux()
	r.Handle("/metrics", promhttp.Handler())
	r.Handle("/alive", &aliveHandler{})
	return r
}
