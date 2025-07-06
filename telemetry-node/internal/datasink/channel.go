package datasink

import (
	"sync"
	"sync/atomic"
)

// Channel is simple wrapper for Go channel
type Channel[T any] struct {
	ch     chan *T
	once   sync.Once
	closed atomic.Bool
}

func NewChannel[T any](bufferSize int) *Channel[T] {
	return &Channel[T]{ch: make(chan *T, bufferSize)}
}

func (c *Channel[T]) Send(value *T) {
	if c.closed.Load() {
		return
	}
	c.ch <- value
}

func (c *Channel[T]) Close() {
	c.once.Do(func() {
		c.closed.Store(true)
		close(c.ch)
	})
}

func (c *Channel[T]) Receive() *T {
	return <-c.ch
}

func (c *Channel[T]) Closed() bool {
	//return c.closed.Load()
	select {
	case <-c.ch:
		return true
	default:
		return false
	}
}
