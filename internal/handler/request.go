package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/binay-das/relay-hub/internal/config"
	"github.com/binay-das/relay-hub/internal/lib"
	"github.com/binay-das/relay-hub/internal/types"
)

func (a App) HandleSendReq(w http.ResponseWriter, r *http.Request) {
	userID, ok := a.userIDFromRequest(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, config.MaxRequestBodySize)
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

	method, rawURL, ok := validateRequestInput(w, payload.Method, payload.URL)
	if !ok {
		return
	}

	if payload.Headers == nil {
		payload.Headers = map[string]string{}
	}

	if payload.Body != "" && (method == "GET" || method == "DELETE") {
		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": method + " requests cannot include a body",
		})
		return
	}

	var body io.Reader

	if payload.Body != "" && (method == "POST" || method == "PUT" || method == "PATCH") {
		body = strings.NewReader(payload.Body)
	}

	req, err := http.NewRequest(method, rawURL, body)

	if err != nil {
		lib.WriteJson(w, http.StatusBadRequest, map[string]any{
			"error":   true,
			"message": "Invalid URL: " + err.Error(),
		})
		return
	}

	for k, v := range payload.Headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{
		Timeout: config.DefaultRequestTimeout,
	}

	start := time.Now()
	response, err := client.Do(req)
	elapsedMs := float64(time.Since(start).Microseconds()) / 1000

	if err != nil {
		_ = a.store.AddHistory(userID, types.RequestHistory{
			Method:      method,
			URL:         rawURL,
			Headers:     payload.Headers,
			RequestBody: payload.Body,
			ElapsedMS:   elapsedMs,
			Error:       true,
			Message:     "Connection error: " + err.Error(),
		})

		lib.WriteJson(w, http.StatusBadGateway, map[string]any{
			"error":      true,
			"message":    "Connection error: " + err.Error(),
			"elapsed_ms": elapsedMs,
		})
		return
	}

	defer response.Body.Close()

	// read response body
	respBytes, err := io.ReadAll(io.LimitReader(response.Body, config.MaxResponseBodySize))

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

	_ = a.store.AddHistory(userID, types.RequestHistory{
		Method:          method,
		URL:             rawURL,
		Headers:         payload.Headers,
		RequestBody:     payload.Body,
		StatusCode:      response.StatusCode,
		StatusText:      http.StatusText(response.StatusCode),
		ElapsedMS:       elapsedMs,
		ResponseHeaders: respHeaders,
		ResponseBody:    bodyStr,
		BodyType:        bodyType,
	})

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
