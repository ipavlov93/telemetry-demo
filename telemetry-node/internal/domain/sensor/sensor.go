package sensor

import "time"

type SamplingRate interface {
	Interval() time.Duration
}
