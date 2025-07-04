package domain_test

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

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

	t.Run("IntervalSensor.Value() happy flow", func(t *testing.T) {
		intervalSensor, err := domain.NewIntervalSensor("", time.Duration(intervalSeconds)*time.Second, randomValue)
		assert.NoError(t, err)

		// context timeout duration has small approximation error or testing code call takes some time
		totalMilliseconds := time.Duration(totalSeconds*1000+1) * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), totalMilliseconds)
		defer cancel()

		resultCh := intervalSensor.Value(ctx)

		var values []*domain.SensorValue
		for value := range resultCh {
			values = append(values, value)
		}

		assert.Equal(t, totalSeconds/intervalSeconds, len(values))
	})
	t.Run("sender should close resultCh after generateFunc panic", func(t *testing.T) {
		intervalSensor, err := domain.NewIntervalSensor("", time.Second, randomPanic)
		assert.NoError(t, err)

		resultCh := intervalSensor.Value(context.Background())

		for range resultCh {
		}
		// no deadlock on resultCh close
	})
}

func TestIntervalSensor_NewIntervalSensor(t *testing.T) {
	t.Run("should return an error when generateFunc is nil", func(t *testing.T) {
		intervalSensor, err := domain.NewIntervalSensor("", time.Second, nil)
		assert.Error(t, err)
		assert.Nil(t, intervalSensor)
	})
}
