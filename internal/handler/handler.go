package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/binay-das/relay-hub/internal/database"
)

const sessionDuration = 30 * 24 * time.Hour

type App struct {
	store database.Store
}

func NewApp(db *sql.DB) App {
	return App{
		store: database.NewStore(db),
	}
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
