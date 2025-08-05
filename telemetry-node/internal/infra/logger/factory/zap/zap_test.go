package zap_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger"
	zapfactory "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/infra/logger/factory/zap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewLogger_WithZapTee(t *testing.T) {
	buff1 := &bytes.Buffer{}
	buff2 := &bytes.Buffer{}

	cfg := logger.ConfigMap{
		"log_destination_1": {
			Enabled:     true,
			MinLevel:    "debug",
			Destination: "buffer1",
		},
		"log_destination_2": {
			Enabled:     true,
			MinLevel:    "error",
			Destination: "buffer2",
		},
	}

	factory := func(out logger.Configuration) (io.Writer, error) {
		switch out.Destination {
		case "buffer1":
			return zapcore.AddSync(buff1), nil
		case "buffer2":
			return zapcore.AddSync(buff2), nil
		default:
			return nil, errors.New("unknown log destination")
		}
	}

	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	log, err := zapfactory.NewLogger(cfg, factory, encoder)
	require.NoError(t, err)

	log.Info("info message")
	log.Error("error message")

	assert.Contains(t, buff1.String(), "info message")
	assert.Contains(t, buff1.String(), "error message")

	assert.NotContains(t, buff2.String(), "info message")
	assert.Contains(t, buff2.String(), "error message")
}
