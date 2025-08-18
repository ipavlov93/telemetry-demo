package builder

import "fmt"

type Strategy interface {
	ConfigBuiltByBuilder() Builder
	Builder
}

type Builder interface {
	ConfigResult() (*Result, error)
	fmt.Stringer
}
