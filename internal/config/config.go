package config

import "time"

var AllowedMethods = map[string]bool{
	"GET":    true,
	"POST":   true,
	"PUT":    true,
	"PATCH":  true,
	"DELETE": true,
}

const (
	MaxRequestBodySize    = 1_000_000 // 1 MB limit for incoming JSON payload
	MaxResponseBodySize   = 5_000_000 // 5 MB limit for reading upstream response body
	DefaultRequestTimeout = 30 * time.Second
	DefaultPort           = ":8080"
)

func IsAllowedMethod(method string) bool {
	return AllowedMethods[method]
}
