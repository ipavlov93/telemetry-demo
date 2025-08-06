package env

import (
	"os"
	"strconv"
	"time"
)

// EnvironmentVariable returns env variable by given key.
// Otherwise, returns fallback.
func EnvironmentVariable(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// ParseBoolEnv lookups and parse env variable by given key.
// Otherwise, returns fallback.
func ParseBoolEnv(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsedValue, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsedValue
}

// ParseIntEnv lookups and parse env variable by given key.
// Otherwise, returns fallback.
func ParseIntEnv(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsedValue, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return fallback
	}
	return int(parsedValue)
}

// ParseUintEnv lookups and parse env variable by given key.
// Otherwise, returns fallback.
func ParseUintEnv(key string, fallback uint) uint {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsedValue, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return fallback
	}
	return uint(parsedValue)
}

// ParseDurationEnv lookups and parse env variable by given key.
// Otherwise, returns fallback.
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

// ParseFloat32Env lookups and parse env variable by given key.
// Otherwise, returns fallback.
func ParseFloat32Env(key string, fallback float32) float32 {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsedValue, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return fallback
	}
	return float32(parsedValue)
}
