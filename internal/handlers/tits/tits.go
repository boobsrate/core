package tits

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	baseHandler
	tits Service
}

func NewTitsHandler(tits Service) *Handler {
	return &Handler{
		tits: tits,
	}
}

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc("/tits", h.listTits).Methods("GET")
	router.HandleFunc("/tits/{cardID}", h.voteTits).Methods("POST")
}

func (h *Handler) listTits(w http.ResponseWriter, r *http.Request) {

	tits, err := h.tits.GetTits(r.Context())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.RespJSON(w, tits, http.StatusOK)
}

func (h *Handler) voteTits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardID"]
	if cardID == "" {
		h.ErrorJSON(w, "", http.StatusBadRequest)
		return
	}

	err := h.tits.IncreaseRating(r.Context(), cardID)
	if err != nil {
		h.ErrorJSON(w, "", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
