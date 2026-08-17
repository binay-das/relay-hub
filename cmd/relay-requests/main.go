package main

import (
	"fmt"
	"net/http"

	"github.com/binay-das/relay-hub/internal/config"
	"github.com/binay-das/relay-hub/internal/handler"
)

func main() {
	http.HandleFunc("/api/request", handler.HandleSendReq)
	http.HandleFunc("/", handler.ServeIndex)

	fmt.Printf("Server running at: http://localhost%s\n", config.DefaultPort)

	err := http.ListenAndServe(config.DefaultPort, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
