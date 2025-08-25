package builder

import (
	"fmt"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/app"
)

// priorityConfigBuilder tries to build app config using given builders with corresponding order
type priorityConfigBuilder struct {
	builders      []app.ConfigBuilder
	resultBuiltBy app.ConfigBuilder
}

func NewPriorityConfigBuilder(builders []app.ConfigBuilder) app.ConfigBuilder {
	return &priorityConfigBuilder{
		builders: builders,
	}
}

// Build sets resultBuiltBy if build has succeeded. Otherwise, it returns an error.
func (b *priorityConfigBuilder) Build() (*app.ConfigResult, error) {
	for _, builder := range b.builders {
		config, err := builder.Build()
		if err != nil {
			continue
		}
		b.resultBuiltBy = builder
		return config, nil
	}
	return nil, fmt.Errorf("no builder found to build app config without error, builder chain length: %d", len(b.builders))
}

func (b *priorityConfigBuilder) BuiltBy() app.ConfigBuilder {
	return b.resultBuiltBy
}

func (b *priorityConfigBuilder) String() string {
	return priorityBuilder.String()
}
