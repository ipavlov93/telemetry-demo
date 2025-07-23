package rate

import (
	"math"
)

// RoundOrDefaultRPS returns the least integer value greater than or equal to given ratePerSecond.
// It will return default value if rps is non-positive.
func RoundOrDefaultRPS(ratePerSecond float64, fallback int) int {
	if ratePerSecond <= 0 {
		return fallback
	}
	return int(math.Ceil(ratePerSecond))
}
