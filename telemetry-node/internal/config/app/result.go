package app

type ConfigResult struct {
	builder   ConfigBuilder
	appConfig *Config
}

func NewConfigResult(builder ConfigBuilder, cfg *Config) *ConfigResult {
	appConfig := cfg
	if appConfig == nil {
		appConfig = &Config{}
	}

	return &ConfigResult{
		builder:   builder,
		appConfig: appConfig,
	}
}

func (r *ConfigResult) AppConfig() *Config { return r.appConfig }
