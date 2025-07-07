// Package retry provides retry strategy components for reliable network communication.
package retry

import (
	"context"
	"time"
)

const defaultRetryAttemptsNumber = 1 // turn off retry attempts by default

// ConstDelayStrategy  retries an operation a fixed retry attempts number with constant delay between each attempt.
type ConstDelayStrategy struct {
	maxAttempts int
	delay       time.Duration
}

// NewConstDelayStrategy returns pointer to created instance of ConstDelayStrategy.
// Important: constructor set defaultRetryAttemptsNumber if maxAttempts is less than 1.
// Notice: constructor skip delay validation or default value set.
func NewConstDelayStrategy(maxAttempts int, delay time.Duration) *ConstDelayStrategy {
	if maxAttempts < 1 {
		maxAttempts = defaultRetryAttemptsNumber
	}
	return &ConstDelayStrategy{
		maxAttempts: maxAttempts,
		delay:       delay,
	}
}

// DoWithContext retries the given operation until:
// - it succeeds;
// - context is done;
// - retry limit is reached.
// It returns the last error if all retries fail.
func (r *ConstDelayStrategy) DoWithContext(parent context.Context, operation func(ctx context.Context) error) error {
	var err error

	for attempt := 1; attempt <= r.maxAttempts; attempt++ {
		// Check if the parent context is already canceled
		select {
		case <-parent.Done():
			return parent.Err()
		default:
		}

		// Do the operation
		err = operation(parent)
		if err == nil {
			return nil // success
		}

		// If not last attempt then wait before next retry attempt
		if attempt < r.maxAttempts {
			select {
			case <-time.After(r.delay):
			case <-parent.Done():
				return parent.Err()
			}
		}
	}

	return err // Return the last error
}
