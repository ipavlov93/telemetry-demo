package builder

import (
	"io"

	jsonpkg "github.com/ipavlov93/telemetry-demo/pkg/json"
)

type configBuilderFromJSON struct {
	reader io.Reader
}

func NewJSONBuilderStrategy(reader io.Reader) Builder {
	return &configBuilderFromJSON{
		reader: reader,
	}
}

func (b *configBuilderFromJSON) ConfigResult() (*Result, error) {
	result := NewResult(b)

	err := jsonpkg.Decode(result.appConfig, b.reader)
	return result, err
}

func (b *configBuilderFromJSON) String() string {
	return builderFromJSON.String()
}
