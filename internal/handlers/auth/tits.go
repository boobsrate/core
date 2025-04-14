package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/boobsrate/core/internal/domain"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type Handler struct {
	baseHandler

	cfKey  string
	isProd bool
}

func NewAuthHandler(centrifugeSignKey, env string) *Handler {
	return &Handler{
		isProd: env == "prod",
		cfKey:  centrifugeSignKey,
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

	jwtData := map[string]interface{}{
		"id":         payload.ID,
		"first_name": payload.FirstName,
		"last_name":  payload.LastName,
		"username":   payload.Username,
		"photo_url":  payload.PhotoUrl,
		"auth_date":  payload.AuthDate,
		"hash":       payload.Hash,
	}

	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour)),
		"iat": jwt.NewNumericDate(time.Now()),
		"sub": jwtData,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(h.cfKey))

	cookieDomain := "dev.boobsrate.com"

	if h.isProd {
		cookieDomain = "boobsrate.com"
	}
	expiration := time.Now().Add(14 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:    "boobs_session",
		Value:   tokenStr,
		Expires: expiration,
		Domain:  cookieDomain,
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(200)
}

func (h *Handler) handleGetToken(w http.ResponseWriter, r *http.Request) {
	// Send token back to frontend

	id := domain.NewID()

	customClaims := jwt.MapClaims{
		"sub": id,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour)),
		"iat": jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	// Sign the JWT using a secret key
	secret := []byte(h.cfKey)
	tokenStr, _ := token.SignedString(secret)

	channel := "boobs_dev"

	if h.isProd {
		channel = "boobs_prod"
	}

	customClaimsChan := jwt.MapClaims{
		"sub":     id,
		"channel": channel,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour)),
		"iat":     jwt.NewNumericDate(time.Now()),
	}

	chanToken := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaimsChan)

	chanTokenStr, _ := chanToken.SignedString(secret)

	chatChannel := "chat_global"

	customClaimsChatChan := jwt.MapClaims{
		"sub":     id,
		"channel": chatChannel,
		"exp":     jwt.NewNumericDate(time.Now().Add(time.Hour)),
		"iat":     jwt.NewNumericDate(time.Now()),
	}

	chatChanToken := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaimsChatChan)

	chatChanTokenStr, _ := chatChanToken.SignedString(secret)

	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr, "chan_token": chanTokenStr, "chat_token": chatChanTokenStr})
}
