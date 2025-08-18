package builder

import "fmt"

type Type string

const (
	builderFromJSON Type = "json_builder"
	builderFromEnv  Type = "env_builder"
	priorityBuilder Type = "priority_builder"
)

var validLogOutputs = map[Type]struct{}{
	builderFromJSON: {},
	builderFromEnv:  {},
	priorityBuilder: {},
}

func (t Type) String() string {
	_, ok := validLogOutputs[t]
	if !ok {
		return fmt.Sprintf("builderType(%s)", string(t))
	}
	return string(t)
}
