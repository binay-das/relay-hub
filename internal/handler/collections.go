package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/binay-das/relay-hub/internal/config"
	"github.com/binay-das/relay-hub/internal/lib"
	"github.com/binay-das/relay-hub/internal/types"
)

func (a App) HandleCollections(w http.ResponseWriter, r *http.Request) {
	userID, ok := a.userIDFromRequest(w, r)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		collections, err := a.store.ListCollections(userID)
		if err != nil {
			lib.WriteServerError(w, err)
			return
		}
		lib.WriteJson(w, http.StatusOK, map[string]any{
			"collections": collections,
		})
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, config.MaxRequestBodySize)
		defer r.Body.Close()

		var payload types.CollectionPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			lib.WriteJson(w, http.StatusBadRequest, map[string]any{
				"error":   true,
				"message": "Invalid JSON request body",
			})
			return
		}

		name := strings.TrimSpace(payload.Name)
		if name == "" {
			lib.WriteJson(w, http.StatusBadRequest, map[string]any{
				"error":   true,
				"message": "Collection name cannot be empty",
			})
			return
		}

		collection, err := a.store.CreateCollection(userID, name)
		if err != nil {
			lib.WriteServerError(w, err)
			return
		}

		lib.WriteJson(w, http.StatusCreated, collection)
	default:
		lib.MethodNotAllowed(w)
	}
}
