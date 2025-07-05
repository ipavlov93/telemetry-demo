package config

import (
	"time"

	"github.com/ipavlov93/telemetry-demo/pkg/env"
)

// Config represents app config
type Config struct {
	SensorName                       string
	SensorInterval                   time.Duration
	GrpcClientMessageBatchPerRequest int64
	GrpcClientMaxRetryAttempts       int64
	GrpcServer                       string
}

// LoadConfigEnv sets Config with environment variables values
func LoadConfigEnv() Config {
	return Config{
		SensorName:                       env.EnvironmentVariable("SENSOR_NAME", "default-sensor-name"),
		SensorInterval:                   env.ParseDurationEnv("SENSOR_INTERVAL", time.Second),
		GrpcClientMessageBatchPerRequest: env.ParseIntEnv("GRPC_CLIENT_MESSAGE_BATCH_PER_REQUEST", 10),
		GrpcClientMaxRetryAttempts:       env.ParseIntEnv("GRPC_CLIENT_MAX_RETRY_ATTEMPT_NUMBER", 5),
		GrpcServer:                       env.EnvironmentVariable("GRPC_SERVER", "localhost:9000"),
	}
}
