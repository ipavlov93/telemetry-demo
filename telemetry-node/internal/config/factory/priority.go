package factory

import (
	"os"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/builder"
)

func NewAppConfigResult(fileName string) (config *builder.Result, err error) {
	// Notice: error check is ignored
	file, err := os.Open(fileName)
	if err == nil {
		defer file.Close()
	}

	builderStrategies := []builder.Builder{
		builder.NewBuilderJSON(file),
		builder.NewBuilderEnv(),
	}

	return builder.NewPriorityBuilder(builderStrategies).Config()
}
