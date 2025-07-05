package env

import (
	"os"
	"strconv"
	"time"
)

// EnvironmentVariable tries to lookup env variable by given key.
// If it's successful, returns string representation of env variable value.
// Otherwise, returns default value (fallback).
func EnvironmentVariable(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// ParseIntEnv tries to lookup env variable by given key.
// If lookup and parse are successful, it returns integer64 representation of env variable value.
// Otherwise, returns default value (fallback).
func ParseIntEnv(key string, fallback int64) int64 {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fallback
	}
	return parsedValue
}

// ParseDurationEnv tries to lookup env variable by given key.
// If lookup and parse are successful, it returns time.Duration representation of env variable value.
// Otherwise, returns default value (fallback).
// Important: valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func ParseDurationEnv(key string, fallback time.Duration) time.Duration {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsedValue, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsedValue
}
