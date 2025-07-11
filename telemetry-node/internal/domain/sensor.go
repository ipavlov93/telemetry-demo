package domain

import (
	"context"
	"sync"
)

// Sensor performs measurements using Run().
type Sensor interface {
	// Run should respect context cancellation (e.g., via <-ctx.Done()) and wait group Done() by design.
	// Return parameters:
	// 1. []SensorValue is sent to channel.
	// 2. Startup errors returned using err (e.g. Run can't start measurements).
	Run(ctx context.Context, wg *sync.WaitGroup) (<-chan []SensorValue, error)
}
