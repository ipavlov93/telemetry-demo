package receive

import (
	"context"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
)

type DrainLastStrategy struct{}

func (d *DrainLastStrategy) Receive(ctx context.Context, valuesChan <-chan []measurement.SensorValue) []measurement.SensorValue {
	for {
		select {
		case <-ctx.Done():
			return nil
		case sensorValues, ok := <-valuesChan:
			if !ok {
				return nil
			}
			return sensorValues
		default:
			return nil
		}
	}
}
