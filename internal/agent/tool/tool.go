package tool

import (
	"context"
)

type Tool interface {
	Name() string
	Description() string
	Parameters() any
	Call(ctx context.Context, input string) (string, error)
}

var Tools = []Tool{
	&Bash{},
}
