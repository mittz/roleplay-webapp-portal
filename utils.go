package main

import (
	"log"
	"os"
	"strconv"
)

func getEnv(key, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}

func getEnvPortalPort() int {
	val := getEnv("PORTAL_PORT", "8080")
	port, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("PORTAL_PORT: %s should be integer: %v", val, err)
	}

	return port
}

func getEnvDataStudioURL() string {
	return getEnv("DATA_STUDIO_URL", "")
}
