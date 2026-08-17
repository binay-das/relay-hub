package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
	http.HandleFunc("/api/request", handleSendReq)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	fmt.Println("Server running at: http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handleSendReq(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1_000_000)
	defer r.Body.Close()

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

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	start := time.Now()
	response, err := client.Do(req)
	elapsedMs := float64(time.Since(start).Microseconds()) / 1000

	if err != nil {
		lib.WriteJson(w, http.StatusBadGateway, map[string]any{
			"error":      true,
			"message":    "Connection error: " + err.Error(),
			"elapsed_ms": elapsedMs,
		})
		return
	}

	defer response.Body.Close()

	// read response body
	respBytes, err := io.ReadAll(io.LimitReader(response.Body, 5_000_000))

	if err != nil {
		lib.WriteJson(w, http.StatusBadGateway, map[string]any{
			"error":   true,
			"message": "Failed to read response body",
		})
		return
	}

	bodyStr := string(respBytes)

	bodyType := "text"

	contentType := response.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		json.Valid(respBytes)
		bodyType = "json"

		var parsed any

		if json.Unmarshal(respBytes, &parsed) == nil {
			if prettyJSON, err := json.MarshalIndent(
				parsed,
				"",
				"  ",
			); err == nil {
				bodyStr = string(prettyJSON)
			}
		}

	}

	respHeaders := make(map[string]string)
	for key, values := range response.Header {
		respHeaders[strings.ToLower(key)] = strings.Join(values, ", ")
	}

	lib.WriteJson(w, http.StatusOK, map[string]any{
		"error":       false,
		"status_code": response.StatusCode,
		"status_text": http.StatusText(response.StatusCode),
		"elapsed_ms":  elapsedMs,
		"headers":     respHeaders,
		"body":        bodyStr,
		"body_type":   bodyType,
	})

}

// func handleIndex(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(
// 		w,
// 		r,
// 		"index.html",
// 	)

// 	if err != nil {
// 		http.Error(
// 			w,
// 			"Failed to render page",
// 			http.StatusInternalServerError,
// 		)

// 		log.Println("Template error:", err)
// 		return
// 	}

// }
