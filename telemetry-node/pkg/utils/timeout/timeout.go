package timeout

import "time"

// TotalTimeout duration for API call
func TotalTimeout(perRetryTimeout time.Duration, totalAttempts uint) time.Duration {
	return perRetryTimeout * time.Duration(totalAttempts)
}
