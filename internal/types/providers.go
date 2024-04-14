package types

import "context"

type Queue[T any] interface {
	Read(ctx context.Context)
	Write(T) error
	Ping() (string, error)
}

type ReadData[T any] struct {
	Data *T
	Err  error
}
