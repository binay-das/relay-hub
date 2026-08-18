package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/binay-das/relay-hub/internal/auth"
	"github.com/binay-das/relay-hub/internal/config"
	"github.com/binay-das/relay-hub/internal/database"
	"github.com/binay-das/relay-hub/internal/types"
)

func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		fmt.Println("JSON encoding errror:", err)
	}
}

func WriteServerError(w http.ResponseWriter, err error) {
	WriteJson(w, http.StatusInternalServerError, map[string]any{
		"error":   true,
		"message": err.Error(),
	})
}

func WriteMaybeNotFound(w http.ResponseWriter, err error) {
	if database.IsNotFound(err) {
		WriteJson(w, http.StatusNotFound, map[string]any{
			"error":   true,
			"message": "Resource not found",
		})
		return
	}

	WriteServerError(w, err)
}

func MethodNotAllowed(w http.ResponseWriter) {
	WriteJson(w, http.StatusMethodNotAllowed, map[string]any{
		"error":   true,
		"message": "Method not allowed",
	})
}

func ReadAuthPayload(w http.ResponseWriter, r *http.Request) (types.AuthPayload, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, config.MaxRequestBodySize)
	defer r.Body.Close()

	var payload types.AuthPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "Invalid JSON request body",
		})
		return types.AuthPayload{}, false
	}

	payload.Email = strings.ToLower(strings.TrimSpace(payload.Email))
	if _, err := mail.ParseAddress(payload.Email); err != nil {
		WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "Enter a valid email",
		})
		return types.AuthPayload{}, false
	}

	if err := auth.ValidatePassword(payload.Password); err != nil {
		WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": err.Error(),
		})
		return types.AuthPayload{}, false
	}

	return payload, true
}
