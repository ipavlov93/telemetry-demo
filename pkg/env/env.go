package env

import (
	"os"
	"strconv"
	"time"
)

// EnvironmentVariable tries to lookup env variable by given key.
// If lookup is successful, it returns parsed value.
// Otherwise, returns default value (fallback).
func EnvironmentVariable(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// ParseIntEnv tries to lookup and parse env variable by given key.
// If lookup and parse are successful, it returns parsed value.
// Otherwise, returns default value (fallback).
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

// ParseUintEnv tries to lookup and parse env variable by given key.
// If lookup and parse are successful, it returns parsed value.
// Otherwise, returns default value (fallback).
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

// ParseDurationEnv tries to lookup and parse env variable by given key.
// If lookup and parse are successful, it returns parsed value.
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

// ParseFloat32Env tries to lookup and parse env variable by given key.
// If lookup and parse are successful, it returns parsed value.
// Otherwise, returns default value (fallback).
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
