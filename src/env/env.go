package env

import (
	"os"
	"strings"
)

// GetEnv returns value of environment variable specified by key or fallback value
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// AllVars return all environment variables
func AllVars() map[string]string {
	vars := make(map[string]string)

	// Fetch all env variables
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		name := variable[0]

		vars[name] = os.Getenv(name)
	}

	return vars
}
