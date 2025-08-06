package retry

import (
	"time"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/pkg/utils/timeout"
)

const perRetryTimeout = 100 * time.Millisecond

// NewTotalTimeout duration for API call including retries and initial attempt.
func NewTotalTimeout(maxRetryAttempts uint) time.Duration {
	totalAttempts := 1 + maxRetryAttempts
	totalTimeout := timeout.TotalTimeout(perRetryTimeout, totalAttempts)

	// Notice: if totalAttempts = 1 or totalTimeout <= perRetryTimeout
	// then retryStrategy will never run

	return totalTimeout
}
