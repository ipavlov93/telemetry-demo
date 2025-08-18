package builder

import (
	"fmt"
)

// priorityBuilderStrategy tries to build app config using given builders with corresponding order
type priorityBuilderStrategy struct {
	builders      []Builder
	resultBuiltBy Builder
}

func NewPriorityBuilderStrategy(builders []Builder) Strategy {
	return &priorityBuilderStrategy{
		builders: builders,
	}
}

func (b *priorityBuilderStrategy) ConfigResult() (*Result, error) {
	for _, builder := range b.builders {
		config, err := builder.ConfigResult()
		if err != nil {
			continue
		}
		b.resultBuiltBy = builder
		return config, nil
	}
	return nil, fmt.Errorf("no builder found to build app config without error, builder chain length: %d", len(b.builders))
}

func (b *priorityBuilderStrategy) String() string {
	return priorityBuilder.String()
}

func (b *priorityBuilderStrategy) ConfigBuiltByBuilder() Builder {
	return b.resultBuiltBy
}
