package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewEncoder() zapcore.Encoder {
	zapConfig := zap.NewProductionEncoderConfig()
	zapConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	return zapcore.NewJSONEncoder(zapConfig)
}
