package handler

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/binay-das/relay-hub/internal/auth"
	"github.com/binay-das/relay-hub/internal/config"
	"github.com/binay-das/relay-hub/internal/lib"
	"github.com/binay-das/relay-hub/internal/types"
)

func (a App) userIDFromRequest(w http.ResponseWriter, r *http.Request) (int64, bool) {
	token, ok := bearerTokenFromRequest(w, r)
	if !ok {
		return 0, false
	}

	userID, err := a.store.UserIDBySession(auth.HashToken(token))
	if err != nil {
		lib.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"error":   true,
			"message": "Invalid or expired session",
		})
		return 0, false
	}

	return userID, true
}

func bearerTokenFromRequest(w http.ResponseWriter, r *http.Request) (string, bool) {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	if header == "" || token == "" || token == header {
		lib.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"error":   true,
			"message": "Missing bearer token",
		})
		return "", false
	}

	return token, true
}

func readAuthPayload(w http.ResponseWriter, r *http.Request) (types.AuthPayload, bool) {
	r.Body = http.MaxBytesReader(w, r.Body, config.MaxRequestBodySize)
	defer r.Body.Close()

	var payload types.AuthPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "Invalid JSON request body",
		})
		return types.AuthPayload{}, false
	}

	payload.Email = strings.ToLower(strings.TrimSpace(payload.Email))
	if _, err := mail.ParseAddress(payload.Email); err != nil {
		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "Enter a valid email",
		})
		return types.AuthPayload{}, false
	}

	if err := auth.ValidatePassword(payload.Password); err != nil {
		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": err.Error(),
		})
		return types.AuthPayload{}, false
	}

	return payload, true
}

func validateRequestInput(w http.ResponseWriter, methodValue string, urlValue string) (string, string, bool) {
	method := strings.ToUpper(strings.TrimSpace(methodValue))

	if !config.IsAllowedMethod(method) {
		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "Unsupported method: " + method,
		})
		return "", "", false
	}

	rawURL := strings.TrimSpace(urlValue)

	if rawURL == "" {
		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "URL cannot be empty",
		})
		return "", "", false
	}

	if !strings.HasPrefix(rawURL, "http://") &&
		!strings.HasPrefix(rawURL, "https://") {
		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "URL must start with http:// or https://",
		})
		return "", "", false
	}

	return method, rawURL, true
}

func parseID(w http.ResponseWriter, raw string) (int64, bool) {
	id, err := strconv.ParseInt(strings.Trim(raw, "/"), 10, 64)
	if err != nil || id <= 0 {
		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "Invalid ID",
		})
		return 0, false
	}

	return id, true
}
