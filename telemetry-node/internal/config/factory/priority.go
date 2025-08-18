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

func newAppConfigResult(fileName string) (config *builder.Result, err error) {
	// Notice: error check is ignored
	file, err := os.Open(fileName)
	if err == nil {
		defer file.Close()
	}

	builders := []builder.Builder{
		builder.NewJSONBuilderStrategy(file),
		builder.NewBuilderEnv("APP_LOGGER"),
	}

	builderStrategy := builder.NewPriorityBuilderStrategy(builders)
	result, err := builderStrategy.ConfigResult()
	if err != nil {
		return nil, err
	}

	if err = result.AppConfig().Validate(); err != nil {
		return nil, fmt.Errorf(
			"builder %v returned invalid app config: %v",
			builderStrategy.ConfigBuiltByBuilder(), err)
	}

	return result, nil
}
