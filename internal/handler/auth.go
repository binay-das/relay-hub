package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/binay-das/relay-hub/internal/auth"
	"github.com/binay-das/relay-hub/internal/lib"
)

func (a App) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		lib.MethodNotAllowed(w)
		return
	}

	payload, ok := readAuthPayload(w, r)
	if !ok {
		return
	}

	passwordHash, err := auth.HashPassword(payload.Password)
	if err != nil {
		lib.WriteServerError(w, err)
		return
	}

	user, err := a.store.CreateUser(payload.Email, passwordHash)
	if err != nil {
		if strings.Contains(err.Error(), "1062") || strings.Contains(err.Error(), "Duplicate") {
			lib.WriteJson(w, http.StatusConflict, map[string]any{
				"error":   true,
				"message": "Email is already registered",
			})
			return
		}
		lib.WriteServerError(w, err)
		return
	}

	token, tokenHash, err := auth.NewSessionToken()
	if err != nil {
		lib.WriteServerError(w, err)
		return
	}
	if err := a.store.CreateSession(user.ID, tokenHash, time.Now().Add(sessionDuration)); err != nil {
		lib.WriteServerError(w, err)
		return
	}

	lib.WriteJson(w, http.StatusCreated, map[string]any{
		"token": token,
		"user":  user,
	})
}

func (a App) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		lib.MethodNotAllowed(w)
		return
	}

	payload, ok := readAuthPayload(w, r)
	if !ok {
		return
	}

	user, passwordHash, err := a.store.GetUserPasswordHash(payload.Email)
	if err != nil || !auth.VerifyPassword(payload.Password, passwordHash) {
		lib.WriteJson(w, http.StatusUnauthorized, map[string]any{
			"error":   true,
			"message": "Invalid email or password",
		})
		return
	}

	token, tokenHash, err := auth.NewSessionToken()
	if err != nil {
		lib.WriteServerError(w, err)
		return
	}
	if err := a.store.CreateSession(user.ID, tokenHash, time.Now().Add(sessionDuration)); err != nil {
		lib.WriteServerError(w, err)
		return
	}

	lib.WriteJson(w, http.StatusOK, map[string]any{
		"token": token,
		"user":  user,
	})
}

func (a App) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		lib.MethodNotAllowed(w)
		return
	}

	token, ok := bearerTokenFromRequest(w, r)
	if !ok {
		return
	}

	if err := a.store.DeleteSession(auth.HashToken(token)); err != nil {
		lib.WriteServerError(w, err)
		return
	}

	lib.WriteJson(w, http.StatusOK, map[string]any{"logged_out": true})
}

func (a App) HandleMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		lib.MethodNotAllowed(w)
		return
	}

	userID, ok := a.userIDFromRequest(w, r)
	if !ok {
		return
	}

	user, err := a.store.GetUserByID(userID)
	if err != nil {
		lib.WriteMaybeNotFound(w, err)
		return
	}

	lib.WriteJson(w, http.StatusOK, map[string]any{"user": user})
}
