package service

import (
	"context"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
)

type DrainStrategy interface {
	// Receive drains channel with respect context done (via <-ctx.Done()) and channel close.
	Receive(ctx context.Context, valuesChan <-chan []measurement.SensorValue) []measurement.SensorValue
}
