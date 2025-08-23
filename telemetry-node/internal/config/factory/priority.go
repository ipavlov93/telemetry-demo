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

	appConfig := configResult.AppConfig()
	if err = appConfig.Validate(); err != nil {
		return nil, fmt.Errorf(
			"builder %v returned invalid app config: %v",
			configResult.BuiltBy(), err)
	}

	return configResult.AppConfig(), nil
}

func newAppConfigResult(fileName string) (config *builder.Result, err error) {
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
