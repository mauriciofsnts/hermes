package worker

import "context"

type Queue[T any] interface {
	Read(ctx context.Context)
	Write(T) error
	Ping() (string, error)
}
