package rate

import (
	"fmt"
	"math/big"
	"time"
)

// samplingRate provides support timeUnits starts with microseconds (nanoseconds are not supported).
// N operations per units (N/units), where rate is numerator and timeUnits is denominator.
type samplingRate struct {
	rate      float32
	timeUnits time.Duration
}

// New constructor will return error if rate or timeUnits are non-positive.
func New(rate float32, timeUnits time.Duration) (*samplingRate, error) {
	if rate <= 0 {
		return nil, fmt.Errorf("can't init samplingRate, rate is invalid")
	}
	if timeUnits <= 0 {
		return nil, fmt.Errorf("can't init samplingRate, timeUnits is invalid")
	}

	return &samplingRate{
		rate:      rate,
		timeUnits: timeUnits,
	}, nil
}

// Interval calculates time interval based on samplingRate.rate
func (r *samplingRate) Interval() time.Duration {
	if !r.valid() {
		return 0
	}

	nanos := big.NewRat(r.timeUnits.Nanoseconds(), 1)
	rate := big.NewRat(1, 1).SetFloat64(float64(r.rate))
	interval := new(big.Rat).Quo(nanos, rate) // interval = nanos / rate

	// Round to nearest integer nanosecond
	nanosecond, _ := interval.Float64()
	return time.Duration(nanosecond)
}

func (r *samplingRate) valid() bool {
	if r == nil {
		return false
	}
	if r.rate <= 0 {
		return false
	}
	return r.timeUnits.Nanoseconds() > 0
}
