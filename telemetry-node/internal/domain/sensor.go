package domain

import "context"

type Sensor interface {
	Value(ctx context.Context) <-chan *SensorValue
}

type DataSink interface {
	Send(data *SensorValue) error
	Close() error
}
