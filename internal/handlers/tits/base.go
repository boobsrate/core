package tits

import (
	"encoding/json"
	"net/http"
)

type baseHandler struct {
}

func (b *baseHandler) ErrorJSON(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if error != "" {
		w.Write([]byte(`{"error":"` + error + `"}`)) // nolint: errcheck
	}
	w.WriteHeader(code)
}

func (b *baseHandler) RespJSON(w http.ResponseWriter, body interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			b.ErrorJSON(w, "", http.StatusInternalServerError)
			return
		}
		w.Write(jsonBody) // nolint: errcheck
	}
}
