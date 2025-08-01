package zap

import (
	"io"

	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func newCore(
	level string,
	defaultLevel zapcore.Level,
	encoder zapcore.Encoder,
	w io.Writer,
) zapcore.Core {
	minLevel := logger.ParseLevelOrDefault(
		level,
		defaultLevel,
	)

	return zapcore.NewCore(
		encoder,
		zapcore.AddSync(w),
		minLevel,
	)
}
