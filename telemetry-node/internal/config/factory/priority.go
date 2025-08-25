package factory

import (
	"fmt"
	"os"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/app"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/builder"
)

func NewAppConfig(fileName string) (config *app.Config, err error) {
	configResult, err := newAppConfigResult(fileName)
	if err != nil {
		return nil, err
	}

	return configResult.AppConfig(), nil
}

func newAppConfigResult(fileName string) (config *app.ConfigResult, err error) {
	// Notice: error check is ignored
	file, err := os.Open(fileName)
	if err == nil {
		defer file.Close()
	}

	builders := []app.ConfigBuilder{
		builder.NewBuilderFromJSON(file),
		builder.NewBuilderFromEnv("APP_LOGGER"),
	}

	priorityConfigBuilder := builder.NewPriorityConfigBuilder(builders)
	result, err := priorityConfigBuilder.Build()
	if err != nil {
		return nil, err
	}

	if err = result.AppConfig().Validate(); err != nil {
		return nil, fmt.Errorf(
			"builder %v returned invalid app config: %v",
			priorityConfigBuilder, err)
	}

	result.SetBuiltBy(priorityConfigBuilder.BuiltBy())

	return result, nil
}
