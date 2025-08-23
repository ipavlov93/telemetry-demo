package builder

import "fmt"

type Type string

const (
	builderTypeJSON Type = "json_builder"
	builderTypeEnv  Type = "env_builder"
)

var validLogOutputs = map[Type]struct{}{
	builderTypeJSON: {},
	builderTypeEnv:  {},
}

func (t Type) String() string {
	_, ok := validLogOutputs[t]
	if !ok {
		return fmt.Sprintf("builderType(%s)", string(t))
	}
	return string(t)
}
