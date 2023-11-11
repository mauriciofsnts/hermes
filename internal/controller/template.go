package controller

import (
	"bytes"
	"encoding/base64"

	tmpl "html/template"

	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/api"
	"github.com/mauriciofsnts/hermes/internal/template"
)

type TemplateController interface {
	Create(c *fiber.Ctx) (string, error)
	Delete(c *fiber.Ctx) error
	ParseTemplate(name string, content map[string]any) (*bytes.Buffer, error)
}

type templateController struct {
	templateService template.TemplateService
}

func NewTemplateController() *templateController {
	return &templateController{
		templateService: template.NewTemplateService(),
	}
}

func (c *templateController) Create(ctx *fiber.Ctx) error {
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

func (c *templateController) GetRaw(ctx *fiber.Ctx) error {
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

func (c *templateController) ParseTemplate(name string, content map[string]any) (*bytes.Buffer, error) {
	html, err := c.templateService.Get(name)

	if err != nil {
		return nil, err
	}

	htmlTmpl, err := tmpl.New(name).Parse(string(html))

	if err != nil {
		return nil, err
	}

	buff := bytes.NewBufferString("")
	err = htmlTmpl.Option("missingkey=error").Execute(buff, content)

	if err != nil {
		return nil, err
	}

	return buff, nil
}
