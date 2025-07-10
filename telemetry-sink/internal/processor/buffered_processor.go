package processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/ipavlov93/telemetry-demo/pkg/grpc/generated/v1/sensor"
)

const (
	defaultBufferSize = 10
)

// BufferedProcessor represent component that save messages to in-memory buffer.
// Flushes messages to outChan when:
// - parent context is done or inputChan is closed (read Run comment);
// - timer tick with given interval;
// - buffer is full (bufferSize is reached).
type BufferedProcessor struct {
	outChan       chan []*pb.SensorValue
	flushInterval time.Duration
	bufferSize    int
}

// NewBufferedProcessor returns pointer to created instance of IntervalSensor.
// Constructor will return error if interval is non-positive.
func NewBufferedProcessor(interval time.Duration, maxSize int) (*BufferedProcessor, error) {
	// required parameters
	if interval <= 0 {
		return nil, fmt.Errorf("can't init BufferedProcessor, interval is invalid")
	}

	// optional parameters
	bufferSize := maxSize
	if maxSize == 0 {
		bufferSize = defaultBufferSize
	}

	return &BufferedProcessor{
		outChan:       make(chan []*pb.SensorValue, 100),
		flushInterval: interval,
		bufferSize:    bufferSize,
	}, nil
}

// Out return the output channel.
func (p *BufferedProcessor) Out() <-chan []*pb.SensorValue { return p.outChan }

// Run starts process messages in a separate goroutine.
// It saves messages to in-memory buffer.
// Flushes messages to outChan when:
// - parent context is done (via <-ctx.Done());
// - inputChan is closed;
// - timer tick with given interval;
// - buffer is full (bufferSize is reached).
// It respects context cancellation (e.g., via <-ctx.Done()) and wait group by design.
func (p *BufferedProcessor) Run(ctx context.Context, inputChan <-chan []*pb.SensorValue, wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
	}

	go func() {
		messageBuffer := make([]*pb.SensorValue, 0, p.bufferSize)

		timer := time.NewTimer(p.flushInterval)
		defer timer.Stop()

		defer func() {
			// gracefully send buffered messages
			if len(messageBuffer) > 0 {
				p.outChan <- messageBuffer
			}

			// ensure that it's single closer
			// receivers will not wait forever on channel close
			close(p.outChan)

			if wg == nil {
				return
			}
			wg.Done()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case messages, ok := <-inputChan:
				if !ok {
					return
				}

				messageBuffer = append(messageBuffer, messages...)
				if len(messageBuffer) >= p.bufferSize {
					p.outChan <- messageBuffer
					// flush buffer and save buffer capacity
					messageBuffer = messageBuffer[:0]

					// to avoid block on Reset if timer already expired or wasn't stopped correctly
					if !timer.Stop() {
						<-timer.C
					}
					timer.Reset(p.flushInterval)
				}
			case <-timer.C:
				if len(messageBuffer) > 0 {
					p.outChan <- messageBuffer
					// flush buffer and save buffer capacity
					messageBuffer = messageBuffer[:0]
				}
				timer.Reset(p.flushInterval)
			}
		}
	}()
}
