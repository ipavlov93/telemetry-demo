package builder

import "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/app"

type Builder interface {
	Config() (*app.Config, error)
}
