package config

import (
	"time"

	"github.com/ipavlov93/telemetry-demo/pkg/env"
)

// Config represents app config
type Config struct {
	SensorName                  string
	SensorInterval              time.Duration
	MessageBatchPerRequest      int64
	GrpcClientMaxRetryAttempts  int64
	GrpcClientRetrySendInterval time.Duration
	GrpcServer                  string
}

// LoadConfigEnv sets Config with environment variables values
func LoadConfigEnv() Config {
	return Config{
		SensorName:                  env.EnvironmentVariable("SENSOR_NAME", "default-sensor-name"),
		SensorInterval:              env.ParseDurationEnv("SENSOR_INTERVAL", time.Second),
		MessageBatchPerRequest:      env.ParseIntEnv("MESSAGE_BATCH_PER_REQUEST", 1),
		GrpcClientMaxRetryAttempts:  env.ParseIntEnv("GRPC_CLIENT_MAX_RETRY_ATTEMPTS_NUMBER", 5),
		GrpcClientRetrySendInterval: env.ParseDurationEnv("GRPC_CLIENT_RETRY_SEND_INTERVAL", time.Second),
		GrpcServer:                  env.EnvironmentVariable("GRPC_SERVER", "localhost:9000"),
	}
}
