package builder

import "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/app"

type Result struct {
	appConfig *app.Config
	builtBy   Builder
}

func NewResult(builder Builder) *Result {
	return &Result{
		builtBy:   builder,
		appConfig: &app.Config{},
	}
}

func NewResultWithConfig(builder Builder, cfg *app.Config) *Result {
	return &Result{
		builtBy:   builder,
		appConfig: cfg,
	}
}

func (r *Result) AppConfig() *app.Config { return r.appConfig }

func (r *Result) BuiltBy() Builder { return r.builtBy }
