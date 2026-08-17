package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/binay-das/relay-hub/internal/lib"
	"github.com/binay-das/relay-hub/internal/types"
)

var allowedMethods = map[string]bool{
	"GET":    true,
	"POST":   true,
	"PUT":    true,
	"PATCH":  true,
	"DELETE": true,
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		res := map[string]string{
			"msg": "Hi! Server is running, try visiting other sites",
		}
		json.NewEncoder(w).Encode(res)
	})

	fmt.Println("Server running at: http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handleSendReq(w http.ResponseWriter, r *http.Request) {
	var payload types.ReqPayLoad

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "Invalid JSON request body",
		})

		return
	}

	method := strings.ToUpper(strings.TrimSpace(payload.Method))

	if !allowedMethods[method] {
		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "Unsupported method: " + method,
		})

		return
	}

	rawURL := strings.TrimSpace(payload.URL)

	if rawURL == "" {
		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "URL cannot be empty",
		})
		return
	}

	if !strings.HasPrefix(rawURL, "http://") &&
		!strings.HasPrefix(rawURL, "https://") {

		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "URL must start with http:// or https://",
		})
		return
	}

	// outgoing req body

	var body io.Reader

	if payload.Body != "" && (method == "POST" || method == "PUT" || method == "PATCH") {
		body = strings.NewReader(payload.Body)
	}

	// create outgoing request

	req, err := http.NewRequest(method, rawURL, body)

	if err != nil {
		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "Invalid URL: " + err.Error(),
		})
		return
	}

	client := &http.Client{}

	_, err = client.Do(req)

	if err != nil {
		lib.WriteJson(w, http.StatusBadGateway, map[string]any{
			"error":   true,
			"message": "Error",
		})
		return
	}
}
