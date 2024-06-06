package api

import "net/http"

type WrappedHandler func(r *http.Request) Response

type Router interface {
	Get(path string, handler WrappedHandler)
	Post(path string, handler WrappedHandler)
	Put(path string, handler WrappedHandler)
	Delete(path string, handler WrappedHandler)
}
