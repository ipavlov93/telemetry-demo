package channel

import (
	"context"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/service"
)

type drainAllStrategy struct{}

func NewDrainAllStrategy() service.DrainStrategy {
	return &drainAllStrategy{}
}

func (d *drainAllStrategy) Receive(ctx context.Context, valuesChan <-chan []measurement.SensorValue) []measurement.SensorValue {
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
