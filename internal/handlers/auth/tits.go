package auth

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	baseHandler
}

func NewAuthHandler() *Handler {
	return &Handler{}
}

type tgUser struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
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
	h.RespJSON(w, jsonBody, http.StatusOK)
}
