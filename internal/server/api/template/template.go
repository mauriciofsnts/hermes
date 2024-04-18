package template

import (
	"bytes"
	"encoding/base64"

	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/providers/template"
	"github.com/mauriciofsnts/hermes/internal/server/helper"
)

type TemplateControllerInterface interface {
	Create(c *fiber.Ctx) (string, error)
	Delete(c *fiber.Ctx) error
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

func (c *TemplateController) Create(ctx *fiber.Ctx) error {
	payload := struct {
		Content string `json:"content"`
		Name    string `json:"name"`
	}{}

	if err := ctx.BodyParser(&payload); err != nil {
		return helper.Err(ctx, fiber.StatusBadRequest, "invalid request body", err)
	}

	parsedContent, err := base64.StdEncoding.DecodeString(payload.Content)

	if err != nil {
		return helper.Err(ctx, fiber.StatusBadRequest, "failed to decode content", err)
	}

	if c.templateProvider.Exists(payload.Name) {
		return helper.Err(ctx, fiber.StatusBadRequest, "template already exists", nil)
	}

	err = c.templateProvider.Create(payload.Name, []byte(parsedContent))

	if err != nil {
		return helper.Err(ctx, fiber.StatusInternalServerError, "failed to create template", err)
	}

	return helper.Success(ctx, fiber.StatusCreated, "template created successfully")
}

func (c *TemplateController) GetRaw(ctx *fiber.Ctx) error {
	id := ctx.Params("slug")

	if id == "" {
		return helper.Err(ctx, fiber.StatusBadRequest, "slug is required", nil)
	}

	html, err := c.templateProvider.Get(id)

	if err != nil {
		return helper.Err(ctx, fiber.StatusInternalServerError, "failed to get template", err)
	}

	ctx.Context().SetContentType("text/html")
	return ctx.Send(html)
}
