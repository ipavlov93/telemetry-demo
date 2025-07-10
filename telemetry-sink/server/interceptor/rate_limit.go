package interceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// ByteRateLimiterInterceptor creates a gRPC interceptor that applies a token bucket rate limit
// based on the size of the incoming request message in bytes.
func ByteRateLimiterInterceptor(limiter *rate.Limiter, lg logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// This is necessary to use proto.Size() to determine the message's marshaled size.
		msg, ok := req.(proto.Message)
		if !ok {
			return nil, fmt.Errorf("request type %T is not a proto.Message", req)
		}

		// Calculate the marshaled size of the protobuf message
		msgSize := proto.Size(msg)

		// If theoretically the message size is 0, it shouldn't consume token
		if msgSize == 0 {
			return handler(ctx, req)
		}

		// Check if the limiter allows the request based on its size.
		if !limiter.AllowN(time.Now(), msgSize) {
			// If the limit is exceeded, log the event and return a gRPC error with codes.ResourceExhausted.
			lg.Debug(fmt.Sprintf("Rate limit exceeded for request of size %d bytes. Limit: %d bytes/sec, Burst: %d bytes.", msgSize, limiter.Limit(), limiter.Burst()))
			return nil, status.Errorf(codes.ResourceExhausted, "Rate limit exceeded. Please reduce your request rate and try again later.")
		}

		// Proceed to the next handler in the chain.
		return handler(ctx, req)
	}
}
