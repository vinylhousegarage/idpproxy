package config

import (
	"os"
)

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return "9000"
	}
	return port
}

func GetOpenAPIURL() string {
	url := os.Getenv("OPENAPI_URL")
	if url == "" {
		panic("OPENAPI_URL is not set")
	}
	return url
}
