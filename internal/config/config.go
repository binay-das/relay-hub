package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

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

type Config struct {
	Env        string `yaml:"env"`
	DBDSN      string `yaml:"db_dsn"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address string `yaml:"address"`
}

func Load() Config {
	data, err := os.ReadFile("config/local.yaml")

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		panic(err)
	}

	return config
}
