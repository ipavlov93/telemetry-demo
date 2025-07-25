package receive

import (
	"context"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
)

// NopDrainStrategy ignores drain strategy
type NopDrainStrategy struct{}

func (n *NopDrainStrategy) Receive(_ context.Context, _ <-chan []measurement.SensorValue) []measurement.SensorValue {
	return nil
}
