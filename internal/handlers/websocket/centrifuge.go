package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func handleGetToken(w http.ResponseWriter, r *http.Request) {
	// Send token back to frontend

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "boobs-backend",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	// Sign the JWT using a secret key
	secret := []byte("UH6zHlXGZcAK6mfYVuVuqe3A5QLq")
	tokenStr, _ := token.SignedString(secret)

	json.NewEncoder(w).Encode(map[string]string{"token": tokenStr})
}
