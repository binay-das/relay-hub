package main

import (
	"fmt"
	"net/http"

	"github.com/binay-das/relay-hub/internal/config"
	"github.com/binay-das/relay-hub/internal/database"
	"github.com/binay-das/relay-hub/internal/handler"
)

func main() {
	cfg := config.Load()

	fmt.Println("Environment:", cfg.Env)
	fmt.Println("Server:", cfg.HTTPServer.Address)
	db, err := database.ConnectToDB(cfg.DBDSN)
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}
	defer db.Close()

	fmt.Println("Database connected!")
	if err := database.Migrate(db); err != nil {
		fmt.Println("Database migration failed:", err)
		return
	}
	fmt.Println("Database migrated!")

	app := handler.NewApp(db)

	http.HandleFunc("/api/auth/register", app.HandleRegister)
	http.HandleFunc("/api/auth/login", app.HandleLogin)
	http.HandleFunc("/api/auth/logout", app.HandleLogout)
	http.HandleFunc("/api/auth/me", app.HandleMe)
	http.HandleFunc("/api/request", app.HandleSendReq)
	http.HandleFunc("/api/collections", app.HandleCollections)
	http.HandleFunc("/api/saved-requests/", app.HandleSavedRequestByID)
	http.HandleFunc("/api/saved-requests", app.HandleSavedRequests)
	http.HandleFunc("/api/history/", app.HandleHistoryByID)
	http.HandleFunc("/api/history", app.HandleHistory)
	http.HandleFunc("/", handler.ServeIndex)

	fmt.Printf("Server running at: http://localhost%s\n", config.DefaultPort)

	err = http.ListenAndServe(config.DefaultPort, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
