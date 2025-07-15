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
