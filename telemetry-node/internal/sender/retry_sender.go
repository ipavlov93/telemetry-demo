// Package sender provides components for reliable network communication.
package sender

import (
	"context"
	"time"
)

const defaultRetryAttemptsNumber = 1 // turn off retry attempts by default

// RetrySender represents component with retry strategy that respects context cancellation by design.
type RetrySender interface {
	DoWithContext(context.Context, func(context.Context) error) error
}

// ConstDelaySender implements the RetrySender interface.
// It retries an operation a fixed retry attempts number with constant delay between each attempt.
type ConstDelaySender struct {
	maxAttempts int
	delay       time.Duration
}

// NewRetrySender returns pointer to created instance of ConstDelaySender.
// Constructor set defaultRetryAttemptsNumber if maxAttempts is less than 1.
// Notice: constructor skip delay validation or default value set.
func NewRetrySender(maxAttempts int, delay time.Duration) *ConstDelaySender {
	if maxAttempts < 1 {
		maxAttempts = defaultRetryAttemptsNumber
	}
	return &ConstDelaySender{
		maxAttempts: maxAttempts,
		delay:       delay,
	}
}

// DoWithContext retries the given operation until:
// - it succeeds;
// - context is done;
// - retry limit is reached.
// It returns the last error if all retries fail.
func (r *ConstDelaySender) DoWithContext(parent context.Context, operation func(ctx context.Context) error) error {
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
