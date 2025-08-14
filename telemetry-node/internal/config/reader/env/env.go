package env

import (
	"github.com/ipavlov93/telemetry-demo/pkg/env"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/app"
)

// LoadOrDefault parses environment variables or set default values.
func LoadOrDefault() *app.Config {
	return &app.Config{
		SensorName:                 env.EnvironmentVariable("SENSOR_NAME", "default-sensor-name"),
		SensorValueRatePerSecond:   env.ParseFloat32Env("SENSOR_VALUE_RATE_PER_SECOND", 5),
		RequestRatePerSecond:       env.ParseFloat32Env("REQUEST_RATE_PER_SECOND", 1),
		GrpcServerSocket:           env.EnvironmentVariable("GRPC_SERVER_SOCKET", "localhost:8000"),
		GrpcClientMaxRetryAttempts: env.ParseUintEnv("GRPC_CLIENT_MAX_RETRY_ATTEMPTS_NUMBER", 1),
	}
}
