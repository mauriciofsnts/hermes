package helper

import (
	"encoding/json"
	"net/http"
)

func Ok[T any](w http.ResponseWriter, detail T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(detail)
}

func Created[T any](w http.ResponseWriter, detail T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(detail)
}

func Err(
	w http.ResponseWriter,
	error ErrorType,
	message string,
) {
	DetailedError(w, error, map[string]string{
		"message": message,
	})
}

func DetailedError[T any](
	w http.ResponseWriter,
	error ErrorType,
	detail T,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(error.StatusCode)
	_ = json.NewEncoder(w).Encode(Error[T]{
		Error:  error.Name,
		Detail: detail,
	})
}
