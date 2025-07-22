package config

import "github.com/ipavlov93/telemetry-demo/pkg/env"

type AppConfig struct {
	SensorName string
	// Number of requests per second to send
	RequestRatePerSecond     float32
	SensorValueRatePerSecond float32
	// telemetry-sink gRPC socket
	GrpcServerSocket           string
	GrpcClientMaxRetryAttempts uint
	LoggerMinLogLevel          string
}

// LoadConfigEnv parses environment variables or set default values.
func LoadConfigEnv() AppConfig {
	return AppConfig{
		SensorName:                 env.EnvironmentVariable("SENSOR_NAME", "default-sensor-name"),
		SensorValueRatePerSecond:   env.ParseFloat32Env("SENSOR_VALUE_RATE_PER_SECOND", 5),
		RequestRatePerSecond:       env.ParseFloat32Env("REQUEST_RATE_PER_SECOND", 1),
		GrpcServerSocket:           env.EnvironmentVariable("GRPC_SERVER_SOCKET", "localhost:8000"),
		GrpcClientMaxRetryAttempts: env.ParseUintEnv("GRPC_CLIENT_MAX_RETRY_ATTEMPTS_NUMBER", 1),
		LoggerMinLogLevel:          env.EnvironmentVariable("LOGGER_MIN_LOG_LEVEL", "info"),
	}
}
