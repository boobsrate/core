package chat

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/boobsrate/core/internal/domain"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2/jws"
)

type Handler struct {
	baseHandler

	cfKey     string
	isProd    bool
	wsChannel chan domain.WSMessage
}

func NewChatHandler(centrifugeSignKey, env string, wsChannel chan domain.WSMessage) *Handler {
	return &Handler{
		isProd:    env == "prod",
		cfKey:     centrifugeSignKey,
		wsChannel: wsChannel,
	}
}

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc("/chat/messages", h.postMessage).Methods("POST")
}

type chatPayload struct {
	Text string `json:"text"`
}

type tgPayload struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	PhotoUrl  string `json:"photo_url"`
	AuthDate  int    `json:"auth_date"`
	Hash      string `json:"hash"`
}

func (h *Handler) postMessage(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		h.ErrorJSON(w, err.Error(), 500)
	}
	var payload chatPayload
	err = json.Unmarshal(jsonBody, &payload)
	if err != nil {
		fmt.Println(err)
		h.ErrorJSON(w, err.Error(), 500)
	}

	// extract cookie boobs_session
	cookie, err := r.Cookie("boobs_session")
	if cookie == nil {
		h.ErrorJSON(w, "no cookie", 500)
	}

	if err != nil {
		fmt.Println(err)
		h.ErrorJSON(w, err.Error(), 500)
	}

	// extract value form cookie
	tokenStrJWT := cookie.Value

	claims, err := jws.Decode(tokenStrJWT)
	if err != nil {
		fmt.Println(err)
		h.ErrorJSON(w, err.Error(), 500)
	}

	// unmarshal claims.Sub to ustgPl
	var ustgPl tgPayload
	err = json.Unmarshal([]byte(claims.Sub), &ustgPl)
	if err != nil {
		fmt.Println(err)
		h.ErrorJSON(w, err.Error(), 500)
	}

	h.wsChannel <- domain.WSMessage{
		Type: domain.WSMessageTypeChat,
		Message: domain.WSChatMessage{
			Text:   payload.Text,
			Sender: ustgPl.Username,
		},
	}

	h.RespJSON(w, map[string]string{"message": "ok"}, http.StatusOK)
}
