package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/binay-das/relay-hub/internal/config"
	"github.com/binay-das/relay-hub/internal/lib"
	"github.com/binay-das/relay-hub/internal/types"
)

func (a App) HandleSavedRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := a.userIDFromRequest(w, r)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		requests, err := a.store.ListSavedRequests(userID)
		if err != nil {
			lib.WriteServerError(w, err)
			return
		}
		lib.WriteJson(w, http.StatusOK, map[string]any{
			"requests": requests,
		})
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, config.MaxRequestBodySize)
		defer r.Body.Close()

		var payload types.SaveRequestPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			lib.WriteJson(w, http.StatusBadRequest, map[string]any{
				"error":   true,
				"message": "Invalid JSON request body",
			})
			return
		}

		method, rawURL, ok := validateRequestInput(w, payload.Method, payload.URL)
		if !ok {
			return
		}

		payload.Method = method
		payload.URL = rawURL
		payload.Name = strings.TrimSpace(payload.Name)
		if payload.Name == "" {
			payload.Name = method + " " + rawURL
		}
		if payload.Headers == nil {
			payload.Headers = map[string]string{}
		}

		request, err := a.store.SaveRequest(userID, payload)
		if err != nil {
			lib.WriteServerError(w, err)
			return
		}

		lib.WriteJson(w, http.StatusCreated, request)
	default:
		lib.MethodNotAllowed(w)
	}
}

func (a App) HandleSavedRequestByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := a.userIDFromRequest(w, r)
	if !ok {
		return
	}

	id, ok := parseID(w, strings.TrimPrefix(r.URL.Path, "/api/saved-requests/"))
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodDelete:
		if err := a.store.DeleteSavedRequest(userID, id); err != nil {
			lib.WriteMaybeNotFound(w, err)
			return
		}
		lib.WriteJson(w, http.StatusOK, map[string]any{"deleted": true})
	default:
		lib.MethodNotAllowed(w)
	}
}
