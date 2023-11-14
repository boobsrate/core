package tits

import (
	"net/http"
	"strconv"

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
	router.HandleFunc("/tits/top/{limit}", h.listTopTits).Methods("GET")
	router.HandleFunc("/tits/{cardID}", h.voteTits).Methods("POST")
	router.HandleFunc("/tits/report/{cardID}", h.reportTits).Methods("POST")
	router.HandleFunc("/tits/abyss/{limit}", h.listAbyssTits).Methods("GET")

}

func (h *Handler) listAbyssTits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rawLimit := vars["limit"]
	if rawLimit == "" {
		h.ErrorJSON(w, "", http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(rawLimit)
	if err != nil {
		h.ErrorJSON(w, "", http.StatusBadRequest)
		return
	}

	if limit < 1 || limit > 100 {
		limit = 100
		return
	}

	tits, err := h.tits.GetTop(r.Context(), limit, true)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.RespJSON(w, tits, http.StatusOK)
}

func (h *Handler) reportTits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cardID := vars["cardID"]
	if cardID == "" {
		h.ErrorJSON(w, "", http.StatusBadRequest)
		return
	}

	err := h.tits.Report(r.Context(), cardID)
	if err != nil {
		h.ErrorJSON(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) listTopTits(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rawLimit := vars["limit"]
	if rawLimit == "" {
		h.ErrorJSON(w, "", http.StatusBadRequest)
		return
	}
	limit, err := strconv.Atoi(rawLimit)
	if err != nil {
		h.ErrorJSON(w, "", http.StatusBadRequest)
		return
	}

	if limit < 1 || limit > 100 {
		limit = 100
		return
	}

	tits, err := h.tits.GetTop(r.Context(), limit, false)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.RespJSON(w, tits, http.StatusOK)
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
