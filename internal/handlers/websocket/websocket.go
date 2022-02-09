package handlers

import (
	"net/http"

	"github.com/boobsrate/core/internal/domain"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

type WebsocketHandler struct {
	log            *otelzap.Logger
	upgrader       websocket.Upgrader
	clientsChannel chan *domain.WSClient
}

func NewWebsocketHandler(log *zap.Logger, clientsChannel chan *domain.WSClient) *WebsocketHandler {
	return &WebsocketHandler{
		log:            otelzap.New(log.Named("websocket_handler")),
		clientsChannel: clientsChannel,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (h *WebsocketHandler) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Ctx(r.Context()).Error("can not upgrade ws connection", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	client := &domain.WSClient{
		ID:         domain.NewID(),
		Connection: conn,
	}
	h.clientsChannel <- client
}

func (h *WebsocketHandler) Register(router *mux.Router) {
	router.HandleFunc("/ws", h.ServeWS).Methods("GET", "CONNECT")
}
