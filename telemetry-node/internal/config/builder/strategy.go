package builder

import "fmt"

type Builder interface {
	Config() (*Result, error)
	fmt.Stringer
}
