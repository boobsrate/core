package websockethub

import (
	"context"
	"sync"
	"time"

	"github.com/boobsrate/core/internal/domain"
	"github.com/centrifugal/centrifuge-go"
	"github.com/gorilla/websocket"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

const (
	pongWait       = 60 * time.Second
	maxMessageSize = 512
)

type WebsocketsHub struct {
	log         *otelzap.Logger
	clientsLock sync.RWMutex
	clients     map[string]*domain.WSClient
	msgChan     chan domain.WSMessage
	clientsChan chan *domain.WSClient
	dead        chan struct{}
	cli         *centrifuge.Client
}

func NewWebsocketsHub(log *zap.Logger) *WebsocketsHub {
	return &WebsocketsHub{
		log:         otelzap.New(log.Named("websockets_hub")),
		clients:     make(map[string]*domain.WSClient),
		msgChan:     make(chan domain.WSMessage),
		clientsChan: make(chan *domain.WSClient),
		dead:        make(chan struct{}),
		cli: centrifuge.NewProtobufClient("ws://centrifuge.centrifuge:8000", centrifuge.Config{
			Token: "UH6zHlXGZcAK6mfYVuVuqe3A5QLq",
		}),
	}
}

func (w *WebsocketsHub) MessagesChannel() chan domain.WSMessage {
	return w.msgChan
}

func (w *WebsocketsHub) ClientsChannel() chan *domain.WSClient {
	return w.clientsChan
}

func (w *WebsocketsHub) Dead() chan struct{} {
	return w.dead
}

func (w *WebsocketsHub) broadcast(msg interface{}) {
	w.log.Info("processing broadcast", zap.Any("msg", msg))
	w.clientsLock.RLock()
	currentClients := make(map[string]*domain.WSClient, len(w.clients))
	for id, client := range w.clients {
		currentClients[id] = client
	}
	w.clientsLock.RUnlock()
	for _, client := range currentClients {
		client.Mu.Lock()
		err := client.Connection.WriteJSON(msg)
		client.Mu.Unlock()
		if err != nil {
			w.log.Error("can not send msg to user", zap.Error(err))
		}
	}
}

func (w *WebsocketsHub) addClient(client *domain.WSClient) {
	w.clientsLock.RLock()
	w.clients[client.ID] = client
	clientsLen := len(w.clients)
	w.clientsLock.RUnlock()
	client.Connection.SetReadLimit(maxMessageSize)
	err := client.Connection.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		w.log.Error("set read deadline", zap.Error(err))
	}
	client.Connection.SetPongHandler(
		func(string) error {
			err = client.Connection.SetReadDeadline(time.Now().Add(pongWait))
			if err != nil {
				w.log.Error("set pong handler", zap.Error(err))
			}
			return err
		},
	)
	go w.reader(client)
	go func() {
		w.msgChan <- domain.WSMessage{
			Type:    domain.WSMessageTypeOnlineUsers,
			Message: domain.WSOnlineUsersMessage{Online: clientsLen},
		}
	}()
}

func (w *WebsocketsHub) removeClient(client *domain.WSClient) {
	w.clientsLock.RLock()
	delete(w.clients, client.ID)
	clientsLen := len(w.clients)
	w.clientsLock.RUnlock()
	go func() {
		w.msgChan <- domain.WSMessage{
			Type:    domain.WSMessageTypeOnlineUsers,
			Message: domain.WSOnlineUsersMessage{Online: clientsLen},
		}
	}()
}

func (w *WebsocketsHub) reader(client *domain.WSClient) {
	defer w.removeClient(client)
	for {
		t, m, err := client.Connection.ReadMessage()
		w.log.Debug("msg read", zap.Any("type", t), zap.Any("msg", m))
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				w.log.Info("unexpected error", zap.Error(err))
			}
			return
		}
	}
}

func (w *WebsocketsHub) processMsg(msg domain.WSMessage) {
	b, _ := msg.MarshalJSON()
	_, err := w.cli.Publish(context.Background(), "boobs_dev", b)
	if err != nil {
		w.log.Error("can not send msg to centrifuge", zap.Error(err))
	}
	w.broadcast(msg)
}

func (w *WebsocketsHub) Run(ctx context.Context) {
	w.log.Info("websocket hub starting...")
	defer close(w.dead)
	defer w.log.Info("websocket hub stopped")
	w.log.Info("websocket hub started")

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case <-ctx.Done():
			return
		case msg := <-w.msgChan:
			go w.processMsg(msg)
		case client := <-w.clientsChan:
			w.addClient(client)
		}
	}
}
