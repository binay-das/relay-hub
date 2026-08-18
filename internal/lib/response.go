package lib

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/binay-das/relay-hub/internal/database"
)

func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println("JSON encoding error:", err)
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
