// Package app contains application config.
package app

import (
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger"
)

type Config struct {
	SensorName string `json:"sensorName"`
	// Number of requests per second to send
	RequestRatePerSecond     float32 `json:"requestRatePerSecond"`
	SensorValueRatePerSecond float32 `json:"sensorValueRatePerSecond"`
	// telemetry-sink gRPC socket
	GrpcServerSocket           string           `json:"grpcServerSocket"`
	GrpcClientMaxRetryAttempts uint             `json:"grpcClientMaxRetryAttempts"`
	AppLoggerConfig            logger.ConfigMap `json:"appLoggerConfig"`
}
