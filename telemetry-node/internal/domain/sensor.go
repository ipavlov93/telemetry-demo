package domain

import "context"

// Sensor performs measurements using StartMeasurement().
// It should respect context cancellation (e.g., via <-ctx.Done()) and stop gracefully (measurement process finish).
type Sensor interface {
	Run(ctx context.Context) error
}

// DataSink is an abstraction of a destination that receives data (e.g. SensorValue records).
type DataSink[T any] interface {
	Send(data *T)
	Close()
}
