package retry

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

// NewInsecure is allowed to use only for local development.
// It returns a DialOptions without configured TLS.
func NewInsecure(options ...grpc.DialOption) []grpc.DialOption {
	return append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

// NewRetryUnaryInterceptor constructs grpc unary client interceptor with retry.
func NewRetryUnaryInterceptor(maxRetryAttempts uint, lg logger.Logger) []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithUnaryInterceptor(
			retry.UnaryClientInterceptor(
				retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
				// WithMax supplies given value minus one (first initial attempt)
				retry.WithMax(maxRetryAttempts+1),
				retry.WithPerRetryTimeout(perRetryTimeout),
				retry.WithOnRetryCallback(func(ctx context.Context, attempt uint, err error) {
					lg.Debug("retry attempt has failed", zap.Uint("retry_attempt", attempt), zap.Error(err))
				}),
			)),
	}
}
