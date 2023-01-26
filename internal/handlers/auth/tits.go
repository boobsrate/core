package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Handler struct {
	baseHandler
}

func NewAuthHandler() *Handler {
	return &Handler{}
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

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc("/auth/tg-login", h.tgLogin).Methods("POST")
}

func (h *Handler) tgLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	jsonBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.ErrorJSON(w, err.Error(), 500)
	}
	var payload tgPayload
	err = json.Unmarshal(jsonBody, &payload)
	if err != nil {
		h.ErrorJSON(w, err.Error(), 500)
	}
	expiration := time.Now().Add(14 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "boobs_session",
		Value:    strconv.Itoa(payload.ID),
		Expires:  expiration,
		Domain:   ".rate-tits.online",
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(200)
}
