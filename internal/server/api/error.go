package api

import "net/http"

type ErrorType struct {
	Name       string `json:"name"`
	StatusCode int    `json:"status_code"`
}

var (
	NotFoundErr       = ErrorType{"NOT_FOUND", http.StatusNotFound}
	InternalServerErr = ErrorType{"INTERNAL_SERVER_ERROR", http.StatusInternalServerError}
	ForbiddenErr      = ErrorType{"FORBIDDEN", http.StatusForbidden}
	BadRequestErr     = ErrorType{"BAD_REQUEST", http.StatusBadRequest}
	ConflitErr        = ErrorType{"CONFLICT", http.StatusConflict}
	UnauthorizedErr   = ErrorType{"UNAUTHORIZED", http.StatusUnauthorized}
	ValidationErr     = ErrorType{"VALIDATION_ERROR", http.StatusUnprocessableEntity}
	TooFastErr        = ErrorType{"TOO_FAST", http.StatusTooManyRequests}
	TooManyReqErr     = ErrorType{"TOO_MANY_REQUESTS", http.StatusTooManyRequests}
)

type Error[T any] struct {
	Error  string `json:"error" example:"NOT_FOUND"`
	Detail T      `json:"detail"`
}
