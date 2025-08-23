package builder

import (
	"io"

	jsonpkg "github.com/ipavlov93/telemetry-demo/pkg/json"
)

type builderJSON struct {
	reader io.Reader
}

func NewBuilderJSON(reader io.Reader) Builder {
	return &builderJSON{
		reader: reader,
	}
}

func (b *builderJSON) Config() (*Result, error) {
	result := NewResult(b)
	err := jsonpkg.Decode(result.appConfig, b.reader)
	return result, err
}

func (b *builderJSON) String() string {
	return builderTypeJSON.String()
}
