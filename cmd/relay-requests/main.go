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

	http.HandleFunc("/api/request", handler.HandleSendReq)
	http.HandleFunc("/", handler.ServeIndex)

	fmt.Printf("Server running at: http://localhost%s\n", config.DefaultPort)

	err = http.ListenAndServe(config.DefaultPort, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
