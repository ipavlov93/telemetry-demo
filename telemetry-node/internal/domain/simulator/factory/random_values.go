package factory

import (
	"math/rand/v2"
	"time"

	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/rate"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/simulator"
	"go.uber.org/zap"
)

const channelCapacity = 100

func NewRandomValueSimulator(
	name string,
	sensorValueRatePerSecond float32,
	lg logger.Logger,
) (*simulator.SensorSimulator, error) {
	samplingRate, err := rate.New(sensorValueRatePerSecond, time.Second)
	if err != nil {
		lg.Fatal("failed to initialize sampling rate", zap.Error(err))
	}

	sensorSimulator, err := simulator.NewWithRate(
		func() int64 { return rand.Int64N(int64(2 << 16)) },
		samplingRate,
		name,
		channelCapacity,
		lg,
	)
	if err != nil {
		return nil, err
	}
	return sensorSimulator, nil
}
