package template

import (
	"bytes"
	"encoding/base64"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mauriciofsnts/hermes/internal/providers/template"
	"github.com/mauriciofsnts/hermes/internal/server/helper"
	"github.com/mauriciofsnts/hermes/internal/server/validator"
)

type TemplateControllerInterface interface {
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	ParseTemplate(name string, content map[string]any) (*bytes.Buffer, error)
}

type TemplateController struct {
	templateProvider template.TemplateProvider
}

func NewTemplateController() *TemplateController {
	return &TemplateController{
		templateProvider: template.NewTemplateService(),
	}
}

type CreateTemplateBody struct {
	Name    string `json:"name" validate:"required"`
	Content string `json:"content" validate:"required"`
}

func (c *TemplateController) Create(w http.ResponseWriter, r *http.Request) {
	body, validationErr := validator.MustGetBody[CreateTemplateBody](r)

	if validationErr != nil {
		helper.DetailedError(w, helper.ValidationErr, validationErr.Details)
		return
	}

	parsedContent, err := base64.StdEncoding.DecodeString(body.Content)

	if err != nil {
		helper.Err(w, helper.BadRequestErr, "Invalid content")
		return
	}

	if c.templateProvider.Exists(body.Name) {
		helper.Err(w, helper.BadRequestErr, "template already exists")
		return
	}

	err = c.templateProvider.Create(body.Name, []byte(parsedContent))

	if err != nil {
		helper.Err(w, helper.InternalServerErr, "failed to create template")
		return
	}

	helper.Created(w, "template created successfully")
}

func (c *TemplateController) GetRaw(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "slug")

	if id == "" {
		helper.Err(w, helper.BadRequestErr, "slug is required")
		return
	}

	html, err := c.templateProvider.Get(id)

	if err != nil {
		helper.Err(w, helper.InternalServerErr, "failed to get template")
		return
	}

	w.Write(html)
}
