package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type Handler struct {
	baseHandler

	cfKey string
}

func NewAuthHandler(centrifugeSignKey string) *Handler {
	return &Handler{
		cfKey: centrifugeSignKey,
	}
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
	router.HandleFunc("/auth/get-token", h.handleGetToken).Methods("GET")

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
		Name:    "boobs_session",
		Value:   strconv.Itoa(payload.ID),
		Expires: expiration,
		Domain:  "dev.rate-tits.online",
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(200)
}

func (h *Handler) handleGetToken(w http.ResponseWriter, r *http.Request) {
	// Send token back to frontend

	customClaims := jwt.MapClaims{
		"channel": "boobs_dev",
		"iss": "boobs-backend",
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour)),
		"iat": jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	// Sign the JWT using a secret key
	secret := []byte(h.cfKey)
	tokenStr, _ := token.SignedString(secret)

	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}
