package components

import (
	"context"
	"sync"
)

type Component interface {
	Start(context.Context, *sync.WaitGroup)
	Name() string
}
