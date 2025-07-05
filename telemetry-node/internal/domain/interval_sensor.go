package domain

import (
	"context"
	"crypto/rand"
	"fmt"
	"reflect"
	"time"
)

// IntervalSensor simulates real sensor that emits value with given interval.
// It sends emitted value to not nil DataSink.
type IntervalSensor struct {
	sensorName   string
	interval     time.Duration
	generateFunc func() int64
	dataSink     DataSink[SensorValue]
}

func NewIntervalSensor(
	name string,
	interval time.Duration,
	generateFunc func() int64,
	dataSink DataSink[SensorValue],
) (*IntervalSensor, error) {
	if generateFunc == nil {
		return nil, fmt.Errorf("can't init IntervalSensor, generateFunc is nil")
	}
	if dataSink == nil || reflect.ValueOf(dataSink).IsNil() {
		return nil, fmt.Errorf("can't init IntervalSensor, DataSink is nil")
	}

	sensorName := name
	if sensorName == "" {
		sensorName = fmt.Sprintf("sensor-%s", rand.Text())
	}

	return &IntervalSensor{
		sensorName:   sensorName,
		interval:     interval,
		generateFunc: generateFunc,
		dataSink:     dataSink,
	}, nil
}

// Run starts producing data at a constant rate in a separate goroutine.
// Method returns error if IntervalSensor hasn't been set properly.
// The measurement process stops when the context is done (via <-ctx.Done())
func (s IntervalSensor) Run(ctx context.Context) error {
	if s.generateFunc == nil {
		return fmt.Errorf("can't generate value, generateFunc is nil")
	}

	if s.dataSink == nil || reflect.ValueOf(s.dataSink).IsNil() {
		return fmt.Errorf("can't send value, dataSink is nil")
	}

	go func(ctx context.Context) {
		defer func() {
			if r := recover(); r != nil {
				// recovered generateFunc() panic
				fmt.Println("Recovered in f", r)
			}
			// ensure consumers/receivers will not wait forever
			s.dataSink.Close()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(s.interval):
				s.dataSink.Send(&SensorValue{
					SensorName: s.sensorName,
					Value:      s.generateFunc(),
					Timestamp:  time.Now(),
				})
			}
		}
	}(ctx)

	return nil
}
