package helpers

import "os"

// getEnv returns the value of the environment variable named by the key.
// If the variable is not present, it returns the `fallback` value.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
