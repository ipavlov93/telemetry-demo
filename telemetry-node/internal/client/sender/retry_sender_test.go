package sender_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/client/sender"
	"github.com/stretchr/testify/assert"
)

func TestRetrySender_SuccessFirstTry(t *testing.T) {
	t.Run("operation should succeed immediately", func(t *testing.T) {
		sender := sender.NewRetrySender(3, 20*time.Millisecond)

		calls := 0
		err := sender.DoWithContext(context.Background(), func(ctx context.Context) error {
			calls++
			return nil // succeed immediately
		})

		assert.NoError(t, err)
		assert.Equal(t, 1, calls)
	})

	t.Run("operation should succeed after retries", func(t *testing.T) {
		sender := sender.NewRetrySender(5, 20*time.Millisecond)

		calls := 0
		err := sender.DoWithContext(context.Background(), func(ctx context.Context) error {
			calls++
			if calls < 3 {
				return errors.New("random error")
			}
			return nil
		})

		assert.NoError(t, err)
		assert.Equal(t, 3, calls)
	})
	t.Run("operation should fail after max retry attempts number reached", func(t *testing.T) {
		sender := sender.NewRetrySender(3, 20*time.Millisecond)

		calls := 0
		expectedErr := errors.New("always fails")
		err := sender.DoWithContext(context.Background(), func(ctx context.Context) error {
			calls++
			return expectedErr
		})

		assert.Equal(t, expectedErr, err)
		assert.Equal(t, 3, calls)
	})
	t.Run("RetrySender should respect context cancellation", func(t *testing.T) {
		sender := sender.NewRetrySender(5, 100*time.Millisecond)

		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		defer cancel()

		calls := 0
		err := sender.DoWithContext(ctx, func(ctx context.Context) error {
			calls++
			return errors.New("random error")
		})

		assert.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Less(t, calls, 5)
	})
}
