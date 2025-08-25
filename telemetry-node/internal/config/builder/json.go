package builder

import (
	"io"

	jsonpkg "github.com/ipavlov93/telemetry-demo/pkg/json"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/app"
)

type configBuilderFromJSON struct {
	reader io.Reader
}

func NewBuilderFromJSON(reader io.Reader) app.ConfigBuilder {
	return &configBuilderFromJSON{
		reader: reader,
	}
}

func (b *configBuilderFromJSON) Build() (*app.ConfigResult, error) {
	appConfig := &app.Config{}
	err := jsonpkg.Decode(appConfig, b.reader)
	if err != nil {
		return nil, err
	}

	return app.NewConfigResult(appConfig)
}

func (b *configBuilderFromJSON) BuiltBy() app.ConfigBuilder { return b }

func (b *configBuilderFromJSON) String() string {
	return builderFromJSON.String()
}
