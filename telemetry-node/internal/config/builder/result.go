package builder

import "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/app"

type Result struct {
	builder   Builder
	appConfig *app.Config
}

func NewResult(builder Builder) *Result {
	return &Result{
		builder:   builder,
		appConfig: &app.Config{},
	}
}

func NewResultWithConfig(builder Builder, cfg *app.Config) *Result {
	return &Result{
		builder:   builder,
		appConfig: cfg,
	}
}

func (r *Result) AppConfig() *app.Config { return r.appConfig }
