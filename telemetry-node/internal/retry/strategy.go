package retry

import "context"

// RetryStrategy represents component with retry strategy that respects context cancellation by design.
type RetryStrategy interface {
	DoWithContext(context.Context, func(context.Context) error) error
}
