package domain

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type WSMessageType string

func (w WSMessageType) String() string {
	return string(w)
}

const (
	WSMessageTypeNewRating = "new_rating"
)

type WSMessage struct {
	Type    WSMessageType
	Message interface{}
}

func (w WSMessage) MarshalJSON() ([]byte, error) {
	data := map[string]interface{}{
		"type":    w.Type.String(),
		"message": w.Message,
	}
	return json.Marshal(data)
}

type WSNewRatingMessage struct {
	TitsID    string `json:"tits_id"`
	NewRating int64  `json:"new_rating"`
}

type WSClient struct {
	ID         string
	Connection *websocket.Conn
}
