package builder

import (
	logfactory "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger/factory"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/reader/env"
)

type configBuilderFromEnv struct {
	loggerConfigMapKey string
}

func NewBuilderEnv(loggerName string) Builder {
	return &configBuilderFromEnv{
		loggerConfigMapKey: loggerName,
	}
}

// ConfigResult builds config Result with default logger.ConfigMap because env does not support logger.Configuration.
func (b *configBuilderFromEnv) ConfigResult() (*Result, error) {
	cfg := env.LoadOrDefault()
	cfg.AppLoggerConfig = logfactory.NewDefaultLoggerConfig(b.loggerConfigMapKey)

	return NewResultWithConfig(b, cfg), nil
}

func (b *configBuilderFromEnv) String() string {
	return builderFromEnv.String()
}
