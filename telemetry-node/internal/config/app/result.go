package app

import "fmt"

type ConfigResult struct {
	appConfig *Config
	builtBy   ConfigBuilder
}

func NewConfigResult(cfg *Config) (*ConfigResult, error) {
	appConfig := cfg
	if appConfig == nil {
		return nil, fmt.Errorf("appConfig is nil")
	}

	return &ConfigResult{
		appConfig: appConfig,
	}, nil
}

func (r *ConfigResult) AppConfig() *Config { return r.appConfig }

func (r *ConfigResult) SetBuiltBy(builder ConfigBuilder) { r.builtBy = builder }
