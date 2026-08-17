package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		fmt.Println("JSON encoding errror:", err)
	}
}
