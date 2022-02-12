package domain

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type WSMessageType string

func (w WSMessageType) String() string {
	return string(w)
}

const (
	WSMessageTypeNewRating   = "new_rating"
	WSMessageTypeOnlineUsers = "online_users"
)

type WSMessage struct {
	Type    WSMessageType
	Message interface{}
}

func (w WSMessage) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"type":       w.Type.String(),
		"message_id": NewID(),
		"message":    w.Message,
	}
	return json.Marshal(data)
}

type WSNewRatingMessage struct {
	TitsID    string `json:"tits_id"`
	NewRating int64  `json:"new_rating"`
}

type WSOnlineUsersMessage struct {
	Online int `json:"online"`
}

type WSClient struct {
	ID         string
	Connection *websocket.Conn
	Mu         sync.Mutex
}
