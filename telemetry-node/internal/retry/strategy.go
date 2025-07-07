package retry

import "context"

// Strategy represents retry strategy that respects context cancellation by design.
type Strategy interface {
	DoWithContext(context.Context, func(context.Context) error) error
}
