package builder

import (
	logfactory "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger/factory"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/reader/env"
)

type builderEnv struct{}

func NewBuilderEnv() Builder { return &builderEnv{} }

func (b *builderEnv) Config() (*Result, error) {
	cfg := env.LoadOrDefault()
	cfg.AppLoggerConfig = logfactory.NewDefaultConfigMap("APP_LOGGER")

	return NewResultWithConfig(b, cfg), nil
}

func (b *builderEnv) String() string {
	return builderTypeEnv.String()
}
