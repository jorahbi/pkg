package op

import "sync"

type ObjectItemInterface interface {
	Reset()
}

func NewPool(fn func() ObjectItemInterface) *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return fn()
		},
	}
}
