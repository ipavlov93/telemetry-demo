package timeout_test

import (
	"testing"
	"time"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/utils/timeout"
	"github.com/stretchr/testify/assert"
)

func TestTotalTimeout(t *testing.T) {
	tests := []struct {
		name                 string
		perRetryTimeout      time.Duration
		maxRetryAttempts     uint
		expectedTotalTimeout time.Duration
	}{
		{
			name:                 "No retries",
			perRetryTimeout:      time.Second,
			maxRetryAttempts:     0,
			expectedTotalTimeout: 1 * time.Second,
		},
		{
			name:                 "Single retry",
			perRetryTimeout:      2 * time.Second,
			maxRetryAttempts:     1,
			expectedTotalTimeout: 4 * time.Second, // 2s * (1 initial + 1 retry)
		},
		{
			name:                 "Multiple retries",
			perRetryTimeout:      500 * time.Millisecond,
			maxRetryAttempts:     3,
			expectedTotalTimeout: 2 * time.Second,
		},
		{
			name:                 "Zero timeout duration",
			perRetryTimeout:      0,
			maxRetryAttempts:     5,
			expectedTotalTimeout: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total := timeout.TotalTimeout(tt.perRetryTimeout, tt.maxRetryAttempts)
			assert.Equal(t, tt.expectedTotalTimeout, total)
		})
	}
}
