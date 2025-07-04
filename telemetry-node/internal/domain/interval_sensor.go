package domain

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"
)

// IntervalSensor simulates real sensor that emits value with given interval.
type IntervalSensor struct {
	sensorName   string
	interval     time.Duration
	generateFunc func() int64
}

func NewIntervalSensor(
	name string,
	interval time.Duration,
	generateFunc func() int64,
) (*IntervalSensor, error) {
	if generateFunc == nil {
		return nil, fmt.Errorf("can't init IntervalSensor, generateFunc is nil")
	}

	sensorName := name
	if sensorName == "" {
		sensorName = fmt.Sprintf("sensor-%s", rand.Text())
	}

	return &IntervalSensor{
		sensorName:   sensorName,
		interval:     interval,
		generateFunc: generateFunc,
	}, nil
}

// Value produces data concurrently with constant rate (until buffered channel is full)
func (s IntervalSensor) Value(ctx context.Context) <-chan *SensorValue {
	if s.generateFunc == nil {
		panic("can't generate value, generateFunc is nil")
	}

	// unbuffered channel is used to prevent immediate block on channel send
	ch := make(chan *SensorValue, 100)

	go func(ctx context.Context, ch chan<- *SensorValue) {
		defer func() {
			if r := recover(); r != nil {
				// recovered generateFunc() panic
			}
			// ensure consumers/receivers will not wait forever
			close(ch)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(s.interval):
				ch <- &SensorValue{
					SensorName: s.sensorName,
					Value:      s.generateFunc(),
					Timestamp:  time.Now(),
				}
			}
		}
	}(ctx, ch)

	return ch
}
