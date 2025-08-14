package builder

import (
	"fmt"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/app"
)

// priorityBuilder tries to build app config using given strategies with corresponding order
type priorityBuilder struct {
	strategies []Builder
}

func NewPriorityBuilder(strategies []Builder) Builder {
	return &priorityBuilder{
		strategies: strategies,
	}
}

func (b priorityBuilder) Config() (*app.Config, error) {
	for _, builderStrategy := range b.strategies {
		config, err := builderStrategy.Config()
		if err != nil {
			continue
		}
		return config, nil
	}
	return nil, fmt.Errorf("no strategy found to build app config")
}
