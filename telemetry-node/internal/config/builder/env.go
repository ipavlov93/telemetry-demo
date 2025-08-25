package builder

import (
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/app"
	logfactory "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger/factory"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/reader/env"
)

type configBuilderFromEnv struct {
	loggerConfigMapKey string
}

func NewBuilderFromEnv(loggerName string) app.ConfigBuilder {
	return &configBuilderFromEnv{
		loggerConfigMapKey: loggerName,
	}
}

// Build ConfigResult with default logger.ConfigMap because env does not support logger.Configuration.
func (b *configBuilderFromEnv) Build() (*app.ConfigResult, error) {
	cfg := env.LoadOrDefault()
	cfg.AppLoggerConfig = logfactory.NewDefaultLoggerConfig(b.loggerConfigMapKey)

	return app.NewConfigResult(cfg)
}

func (b *configBuilderFromEnv) BuiltBy() app.ConfigBuilder { return b }

func (b *configBuilderFromEnv) String() string {
	return builderFromEnv.String()
}
