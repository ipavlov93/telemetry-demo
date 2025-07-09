package config

import (
	"github.com/ipavlov93/telemetry-demo/pkg/env"
)

// Config represents app config
type Config struct {
	SensorName string
	// Number of requests per second to send
	RequestRatePerSecond float32
	// telemetry-sink gRPC socket
	GrpcServerSocket           string
	GrpcClientMaxRetryAttempts uint
	LoggerMinLogLevel          string
}

// LoadConfigEnv sets Config with environment variables values
func LoadConfigEnv() Config {
	return Config{
		SensorName:                 env.EnvironmentVariable("SENSOR_NAME", "default-sensor-name"),
		RequestRatePerSecond:       env.ParseFloat32Env("REQUEST_RATE_PER_SECOND", 1),
		GrpcServerSocket:           env.EnvironmentVariable("GRPC_SERVER_SOCKET", "localhost:8000"),
		GrpcClientMaxRetryAttempts: env.ParseUintEnv("GRPC_CLIENT_MAX_RETRY_ATTEMPTS_NUMBER", 5),
		LoggerMinLogLevel:          env.EnvironmentVariable("LOGGER_MIN_LOG_LEVEL", "info"),
	}
}
