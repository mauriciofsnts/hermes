package validator

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mauriciofsnts/hermes/internal/server/helper"
)

type APIBodyValidationError struct {
	Details any
	Error   helper.ErrorType
}

func MustGetBody[T any](r *http.Request) (T, *APIBodyValidationError) {
	var body T

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		slog.Error("Failed to decode body", "err", err)
		return body, &APIBodyValidationError{
			Error:   helper.BadRequestErr,
			Details: map[string]string{"message": err.Error()},
		}
	}

	validationErrors := Validate(body)

	if len(validationErrors) > 0 {
		return body, &APIBodyValidationError{
			Error:   helper.ValidationErr,
			Details: validationErrors,
		}
	}

	return body, nil
}
