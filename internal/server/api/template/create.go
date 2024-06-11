package template

import (
	"encoding/base64"
	"net/http"

	"github.com/mauriciofsnts/hermes/internal/server/api"
	"github.com/mauriciofsnts/hermes/internal/server/validator"
)

type CreateTemplateBody struct {
	Name    string `json:"name" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (t *TemplateController) CreateTemplate(r *http.Request) api.Response {
	body, validationErr := validator.MustGetBody[CreateTemplateBody](r)

	if validationErr != nil {
		return api.DetailedError(validationErr.Error, validationErr.Details)
	}

	parsedContent, err := base64.StdEncoding.DecodeString(body.Content)

	if err != nil {
		return api.Err(api.BadRequestErr, "Error on decoding content")
	}

	if t.provider.Exists(body.Name) {
		return api.Err(api.BadRequestErr, "An template with this name already exists")
	}

	err = t.provider.Create(body.Name, []byte(parsedContent))

	if err != nil {
		return api.Err(api.InternalServerErr, "Failed to create template")
	}

	return api.Created("Template created successfully")
}
