package rate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInterval(t *testing.T) {
	tests := []struct {
		name    string
		rate    *samplingRate
		want    time.Duration
		wantErr bool
	}{
		{
			name: "should return value",
			rate: &samplingRate{
				rate:      1,
				timeUnits: time.Nanosecond,
			},
			want: time.Nanosecond,
		},
		{
			name: "should return zero interval",
			rate: &samplingRate{
				rate:      2,
				timeUnits: time.Nanosecond,
			},
			want: 0,
		},
		{
			name: "should return zero interval",
			rate: &samplingRate{
				rate:      1000,
				timeUnits: time.Nanosecond,
			},
			want: 0,
		},
		{
			name: "should return zero interval",
			rate: &samplingRate{
				rate:      10e10,
				timeUnits: time.Nanosecond,
			},
			want: 0,
		},
		{
			name: "should return value",
			rate: &samplingRate{
				rate:      1000,
				timeUnits: time.Microsecond,
			},
			want: time.Nanosecond,
		},
		{
			name: "should return value",
			rate: &samplingRate{
				rate:      1000,
				timeUnits: time.Millisecond,
			},
			want: time.Microsecond,
		},
		{
			name: "should return value",
			rate: &samplingRate{
				rate:      1000,
				timeUnits: time.Second,
			},
			want: time.Millisecond,
		},
		{
			name: "should return value",
			rate: &samplingRate{
				rate:      60,
				timeUnits: time.Minute,
			},
			want: time.Second,
		},
		{
			name: "should return value",
			rate: &samplingRate{
				rate:      60,
				timeUnits: time.Hour,
			},
			want: time.Minute,
		},
		{
			name: "should return zero",
			rate: nil,
			want: 0,
		},
		{
			name: "should return zero if given samplingRate has non-positive field values",
			rate: &samplingRate{
				rate:      0,
				timeUnits: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInterval := tt.rate.Interval()
			assert.EqualValues(t, tt.want, gotInterval)
		})
	}
}
