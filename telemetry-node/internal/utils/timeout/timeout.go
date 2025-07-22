package timeout

import "time"

// TotalTimeout duration for the gRPC call including retries and initial attempt.
func TotalTimeout(perRetryTimeout time.Duration, maxRetryAttemptsNumber uint) time.Duration {
	totalAttempts := 1 + time.Duration(maxRetryAttemptsNumber)
	return perRetryTimeout * totalAttempts
}
