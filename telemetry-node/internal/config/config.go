package config

import (
	"github.com/ipavlov93/telemetry-demo/pkg/env"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger"
)

type AppConfig struct {
	SensorName string
	// Number of requests per second to send
	RequestRatePerSecond     float32
	SensorValueRatePerSecond float32
	// telemetry-sink gRPC socket
	GrpcServerSocket           string
	GrpcClientMaxRetryAttempts uint
	AppLoggerConfig            logger.Config
}

func NewAppConfig() *AppConfig {
	return loadFromEnvVariables()
}

// loadConfigEnv parses environment variables or set default values.
func loadFromEnvVariables() *AppConfig {
	return &AppConfig{
		SensorName:                 env.EnvironmentVariable("SENSOR_NAME", "default-sensor-name"),
		SensorValueRatePerSecond:   env.ParseFloat32Env("SENSOR_VALUE_RATE_PER_SECOND", 5),
		RequestRatePerSecond:       env.ParseFloat32Env("REQUEST_RATE_PER_SECOND", 1),
		GrpcServerSocket:           env.EnvironmentVariable("GRPC_SERVER_SOCKET", "localhost:8000"),
		GrpcClientMaxRetryAttempts: env.ParseUintEnv("GRPC_CLIENT_MAX_RETRY_ATTEMPTS_NUMBER", 1),
		AppLoggerConfig: logger.Config{
			MinLevel: env.EnvironmentVariable("LOGGER_MIN_LOG_LEVEL", "info"),
		},
	}
}
