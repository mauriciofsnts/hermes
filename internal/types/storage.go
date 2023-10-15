package types

type Storage[T any] interface {
	Read()
	Write(T) error
}
