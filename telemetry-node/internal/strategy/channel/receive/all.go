package receive

import (
	"context"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
)

type DrainAllStrategy struct{}

func (d *DrainAllStrategy) Receive(ctx context.Context, valuesChan <-chan []measurement.SensorValue) []measurement.SensorValue {
	var valuesBatch []measurement.SensorValue

	for {
		select {
		case <-ctx.Done():
			return valuesBatch
		case sensorValues, ok := <-valuesChan:
			if !ok {
				return nil
			}
			valuesBatch = append(valuesBatch, sensorValues...)
		default:
			return nil
		}
	}
}
