package channel

import (
	"context"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/service"
)

// nopDrainStrategy ignores drain strategy
type nopDrainStrategy struct{}

func NewNopDrainStrategy() service.DrainStrategy {
	return &nopDrainStrategy{}
}

func (n *nopDrainStrategy) Receive(_ context.Context, _ <-chan []measurement.SensorValue) []measurement.SensorValue {
	return nil
}
