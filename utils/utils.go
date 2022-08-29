package utils

import (
	crand "crypto/rand"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"
)

func getEnv(key, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}

func GetEnvPortalPort() int {
	val := getEnv("PORTAL_PORT", "8080")
	port, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("PORTAL_PORT: %s should be integer: %v", val, err)
	}

	return port
}

func GetEnvDataStudioURL() string {
	return getEnv("DATA_STUDIO_URL", "")
}

func GetEnvProjectID() string {
	projectID := getEnv("PROJECT_ID", "")
	if projectID == "" {
		log.Fatalf("projectID should be set")
	}

	return projectID
}

func GetEnvInstanceConnectionName() string {
	return getEnv("INSTANCE_CONNECTION_NAME", "")
}

func GetEnvDatabaseName() string {
	return getEnv("DATABASE_NAME", "")
}

func GetEnvDatabaseUser() string {
	return getEnv("DATABASE_USER", "")
}

func RandomString(n int) string {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	s := make([]rune, n)

	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}
