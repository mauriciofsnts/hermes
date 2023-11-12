package types

import "context"

type Storage[T any] interface {
	Read(ctx context.Context)
	Write(T) error
	Ping() (string, error)
}
