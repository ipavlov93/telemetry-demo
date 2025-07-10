package config

import (
	"time"

	"github.com/ipavlov93/telemetry-demo/pkg/env"
)

// Config represents app config
type Config struct {
	FilePath            string // Path to the output log file
	BufferSize          int    // Buffer size in bytes
	BufferFlushInterval time.Duration
	GrpcServerSocket    string
	RatePerSecond       int // Number of bytes per second to receive
	LoggerMinLogLevel   string
}

// LoadConfigEnv sets Config with environment variables values
func LoadConfigEnv() Config {
	return Config{
		FilePath:            env.EnvironmentVariable("FILE_PATH", "."),
		BufferSize:          env.ParseIntEnv("BUFFER_SIZE", 10),
		BufferFlushInterval: env.ParseDurationEnv("BUFFER_FLUSH_INTERVAL", time.Minute),
		GrpcServerSocket:    env.EnvironmentVariable("GRPC_SERVER_SOCKET", "localhost:8000"),
		RatePerSecond:       env.ParseIntEnv("RATE_PER_SECOND", 1000),
		LoggerMinLogLevel:   env.EnvironmentVariable("LOGGER_MIN_LOG_LEVEL", "info"),
	}
}
