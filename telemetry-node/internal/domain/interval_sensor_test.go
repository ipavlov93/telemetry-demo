package domain_test

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/datasink"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain"
	"github.com/stretchr/testify/assert"
)

func randomValue() int64 {
	return rand.Int64N(int64(2 << 16))
}

func randomPanic() int64 {
	panic("any reason")
}

func TestIntervalSensor_Value(t *testing.T) {
	intervalSeconds := 1
	totalSeconds := 5

	t.Run("IntervalSensor.Run() happy flow", func(t *testing.T) {
		channelDataSink := datasink.NewChannel[domain.SensorValue](100)

		intervalSensor, err := domain.NewIntervalSensor(
			"",
			time.Duration(intervalSeconds)*time.Second,
			randomValue,
			channelDataSink,
		)
		assert.NoError(t, err)

		// context timeout duration has small approximation error or testing code call takes some time
		totalMilliseconds := time.Duration(totalSeconds*1000+1) * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), totalMilliseconds)
		defer cancel()

		err = intervalSensor.Run(ctx)
		assert.NoError(t, err)

		var values []*domain.SensorValue
		for i := 0; i < totalSeconds; i++ {
			values = append(values, channelDataSink.Receive())
		}

		assert.Equal(t, totalSeconds/intervalSeconds, len(values))
	})
	t.Run("sender should close resultCh after generateFunc panic", func(t *testing.T) {
		channelDataSink := datasink.NewChannel[domain.SensorValue](100)

		intervalSensor, err := domain.NewIntervalSensor(
			"",
			time.Second,
			randomPanic,
			channelDataSink,
		)
		assert.NoError(t, err)

		err = intervalSensor.Run(context.Background())
		assert.NoError(t, err)

		assert.True(t, channelDataSink.Closed())
	})
}

func TestIntervalSensor_NewIntervalSensor_Dummy(t *testing.T) {
	t.Run("should return an error when generateFunc is nil", func(t *testing.T) {
		intervalSensor, err := domain.NewIntervalSensor(
			"",
			time.Second,
			nil,
			nil,
		)
		assert.Error(t, err)
		assert.Nil(t, intervalSensor)
	})
	t.Run("should return an error when dataSink is nil", func(t *testing.T) {
		intervalSensor, err := domain.NewIntervalSensor(
			"",
			time.Second,
			randomValue,
			nil,
		)
		assert.Error(t, err)
		assert.Nil(t, intervalSensor)
	})
}
