package service

import (
	"context"
	"sync"
	"time"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
)

type RateLimiter interface {
	Wait(ctx context.Context) error
}

type RunConfig struct {
	valuesChan      <-chan []measurement.SensorValue
	totalTimeoutRPC time.Duration
	limiter         RateLimiter
	wg              *sync.WaitGroup
}

func NewRunConfig(
	ch <-chan []measurement.SensorValue,
	totalTimeoutRPC time.Duration,
	limiter RateLimiter,
	wg *sync.WaitGroup,
) *RunConfig {
	return &RunConfig{
		valuesChan:      ch,
		totalTimeoutRPC: totalTimeoutRPC,
		limiter:         limiter,
		wg:              wg,
	}
}

func (r *RunConfig) Valid() bool {
	if r == nil {
		return false
	}
	if r.valuesChan == nil {
		return false
	}
	if r.totalTimeoutRPC <= 0 {
		return false
	}
	return r.wg != nil
}
