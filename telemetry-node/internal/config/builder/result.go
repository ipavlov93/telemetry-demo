package builder

import "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/app"

type Result struct {
	Builder
	AppConfig *app.Config
}

func NewResult(builder Builder) *Result {
	return &Result{
		Builder:   builder,
		AppConfig: &app.Config{},
	}
}

func NewResultWithConfig(builder Builder, cfg *app.Config) *Result {
	return &Result{
		Builder:   builder,
		AppConfig: cfg,
	}
}
