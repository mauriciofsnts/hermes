package controller

import (
	"bytes"
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/api"
	"github.com/mauriciofsnts/hermes/internal/template"
)

type TemplateControllerInterface interface {
	Create(c *fiber.Ctx) (string, error)
	Delete(c *fiber.Ctx) error
	ParseTemplate(name string, content map[string]any) (*bytes.Buffer, error)
}

type TemplateController struct {
	templateService template.TemplateServiceInterface
}

func NewTemplateController() *TemplateController {
	return &TemplateController{
		templateService: template.NewTemplateService(),
	}
}

func (c *TemplateController) Create(ctx *fiber.Ctx) error {
	payload := struct {
		Content string `json:"content"`
		Name    string `json:"name"`
	}{}

	if err := ctx.BodyParser(&payload); err != nil {
		return api.Err(ctx, fiber.StatusBadRequest, "invalid request body", err)
	}

	parsedContent, err := base64.StdEncoding.DecodeString(payload.Content)

	if err != nil {
		return api.Err(ctx, fiber.StatusBadRequest, "failed to decode content", err)
	}

	if c.templateService.Exists(payload.Name) {
		return api.Err(ctx, fiber.StatusBadRequest, "template already exists", nil)
	}

	err = c.templateService.Create(payload.Name, []byte(parsedContent))

	if err != nil {
		return api.Err(ctx, fiber.StatusInternalServerError, "failed to create template", err)
	}

	return api.Success(ctx, fiber.StatusCreated, "template created successfully")
}

func (c *TemplateController) GetRaw(ctx *fiber.Ctx) error {
	id := ctx.Params("slug")

	if id == "" {
		return api.Err(ctx, fiber.StatusBadRequest, "slug is required", nil)
	}

	html, err := c.templateService.Get(id)

	if err != nil {
		return api.Err(ctx, fiber.StatusInternalServerError, "failed to get template", err)
	}

	ctx.Context().SetContentType("text/html")
	return ctx.Send(html)
}
