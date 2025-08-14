package builder

import (
	"io"

	jsonpkg "github.com/ipavlov93/telemetry-demo/pkg/json"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/app"
)

type builderJSON struct {
	reader io.Reader
}

func NewBuilderJSON(reader io.Reader) Builder {
	return &builderJSON{
		reader: reader,
	}
}

func (b builderJSON) Config() (*app.Config, error) {
	cfg := &app.Config{}
	return cfg, jsonpkg.Decode(cfg, b.reader)
}
