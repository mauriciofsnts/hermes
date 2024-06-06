package validator

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mauriciofsnts/hermes/internal/server/api"
)

type APIBodyValidationError struct {
	Details any
	Error   api.ErrorType
}

func MustGetBody[T any](r *http.Request) (T, *APIBodyValidationError) {
	var body T

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		slog.Error("Failed to decode body", "err", err)
		return body, &APIBodyValidationError{
			Error:   api.BadRequestErr,
			Details: map[string]string{"message": err.Error()},
		}
	}

	validationErrors := Validate(body)

	if len(validationErrors) > 0 {
		return body, &APIBodyValidationError{
			Error:   api.ValidationErr,
			Details: validationErrors,
		}
	}

	return body, nil
}
