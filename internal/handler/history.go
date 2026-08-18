package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/binay-das/relay-hub/internal/lib"
)

func (a App) HandleHistory(w http.ResponseWriter, r *http.Request) {
	userID, ok := a.userIDFromRequest(w, r)
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodGet:
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		history, err := a.store.ListHistory(userID, limit)
		if err != nil {
			lib.WriteServerError(w, err)
			return
		}
		lib.WriteJson(w, http.StatusOK, map[string]any{
			"history": history,
		})
	case http.MethodDelete:
		if err := a.store.ClearHistory(userID); err != nil {
			lib.WriteServerError(w, err)
			return
		}
		lib.WriteJson(w, http.StatusOK, map[string]any{"deleted": true})
	default:
		lib.MethodNotAllowed(w)
	}
}

func (a App) HandleHistoryByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := a.userIDFromRequest(w, r)
	if !ok {
		return
	}

	id, ok := parseID(w, strings.TrimPrefix(r.URL.Path, "/api/history/"))
	if !ok {
		return
	}

	switch r.Method {
	case http.MethodDelete:
		if err := a.store.DeleteHistory(userID, id); err != nil {
			lib.WriteMaybeNotFound(w, err)
			return
		}
		lib.WriteJson(w, http.StatusOK, map[string]any{"deleted": true})
	default:
		lib.MethodNotAllowed(w)
	}
}
