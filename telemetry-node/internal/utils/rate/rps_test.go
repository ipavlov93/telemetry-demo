package rate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundRPS(t *testing.T) {
	tests := []struct {
		name       string
		rps        float64
		defaultRPS int
		wantRPS    float64
	}{
		{
			name:       "should ceil to 1",
			rps:        0.1,
			defaultRPS: 5,
			wantRPS:    1,
		},
		{
			name:       "should left 1 unchanged",
			rps:        1,
			defaultRPS: 5,
			wantRPS:    1,
		},
		{
			name:       "should set default rps",
			rps:        -5.5,
			defaultRPS: 1,
			wantRPS:    1,
		},
		{
			name:       "should set default rps",
			rps:        0,
			defaultRPS: 5,
			wantRPS:    5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRPS := RoundOrDefaultRPS(tt.rps, tt.defaultRPS)
			assert.EqualValues(t, tt.wantRPS, gotRPS)
		})
	}
}
