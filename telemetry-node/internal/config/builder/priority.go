package builder

import (
	"fmt"
)

// priorityBuilder tries to build app config using given strategies with corresponding order
type priorityBuilder struct {
	strategies    []Builder
	resultBuiltBy Builder
}

func NewPriorityBuilder(strategies []Builder) Builder {
	return &priorityBuilder{
		strategies: strategies,
	}
}

func (b *priorityBuilder) Config() (*Result, error) {
	for _, builderStrategy := range b.strategies {
		config, err := builderStrategy.Config()
		if err != nil {
			continue
		}
		b.resultBuiltBy = builderStrategy
		return config, nil
	}
	return nil, fmt.Errorf("no strategy found to build app config")
}

func (b *priorityBuilder) String() string {
	return b.resultBuiltBy.String()
}
