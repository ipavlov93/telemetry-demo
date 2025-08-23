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
	result := app.NewConfigResult(b, nil)

	err := jsonpkg.Decode(result.AppConfig(), b.reader)
	return result, err
}

func (b *configBuilderFromJSON) String() string {
	return builderFromJSON.String()
}
