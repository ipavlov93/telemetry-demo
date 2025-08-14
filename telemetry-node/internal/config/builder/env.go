package builder

import (
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/app"
	logfactory "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger/factory"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/reader/env"
)

type builderEnv struct{}

func NewBuilderEnv() Builder { return &builderEnv{} }

func (b builderEnv) Config() (*app.Config, error) {
	cfg := env.LoadOrDefault()
	cfg.AppLoggerConfig = logfactory.NewDefaultConfigMap("APP_LOGGER")

	return cfg, nil
}
