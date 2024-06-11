package api

import "net/http"

type Response struct {
	StatusCode int
	Header     http.Header
	Body       any
}
