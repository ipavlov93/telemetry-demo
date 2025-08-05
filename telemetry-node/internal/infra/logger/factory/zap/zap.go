package zap

import (
	"fmt"

	config "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/infra/logger/factory/writer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const defaultMinLogLevel = zapcore.InfoLevel

// NewLogger factory constructs logger with multiple log destinations with corresponding log level.
func NewLogger(
	configMap config.ConfigMap,
	writerFactory writer.Factory,
	encoder zapcore.Encoder,
	option ...zap.Option,
) (*zap.Logger, error) {
	var cores []zapcore.Core

	for _, logCfg := range configMap {
		if !logCfg.Enabled {
			continue
		}

		w, err := writerFactory(logCfg)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to create writer for given destination %q: %v", logCfg.Destination, err)
		}

		cores = append(cores, newCore(logCfg.MinLevel, defaultMinLogLevel, encoder, w))
	}

	if len(cores) == 0 {
		return zap.NewNop(), nil
	}

	tee := zapcore.NewTee(cores...)
	return zap.New(tee, option...), nil
}
