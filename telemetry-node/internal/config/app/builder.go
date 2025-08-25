package app

type ConfigBuilder interface {
	Build() (*ConfigResult, error)
	BuiltBy() ConfigBuilder
}
