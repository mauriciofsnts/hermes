package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/mauriciofsnts/hermes/internal/api"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/queue"
	"github.com/mauriciofsnts/hermes/internal/template"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type EmailControllerInterface interface {
	SendPlainTextEmail(c *fiber.Ctx) error
	SendTemplateEmail(c *fiber.Ctx) error
}

type EmailControler struct {
	EmailControllerInterface
	validate        *validator.Validate
	templateService template.TemplateServiceInterface
}

func NewEmailController() *EmailControler {
	return &EmailControler{
		validate:        validator.New(),
		templateService: template.NewTemplateService(),
	}
}

func (e *EmailControler) SendPlainTextEmail(c *fiber.Ctx) error {
	var bodyEmail types.PlainTextEmail

	if err := c.BodyParser(&bodyEmail); err != nil {
		return api.Err(c, fiber.StatusBadRequest, "Invalid body", err)
	}

	if err := e.validate.Struct(bodyEmail); err != nil {
		return api.Err(c, fiber.StatusBadRequest, "Invalid body", err)
	}

	mail := types.Mail{
		To:      []string{bodyEmail.To},
		Subject: bodyEmail.Subject,
		Sender:  config.Envs.DefaultFrom,
		Body:    bodyEmail.Body,
		Type:    types.TEXT,
	}

	queue.Queue.Write(mail)

	api.Success(c, fiber.StatusCreated, "Email sent successfully")
	return nil
}

func (e *EmailControler) SendTemplateEmail(ctx *fiber.Ctx) error {
	var templateEmail types.TemplateEmail

	templateName := ctx.Params("slug")

	if templateName == "" {
		return api.Err(ctx, fiber.StatusBadRequest, "Invalid template name", nil)
	}

	if err := ctx.BodyParser(&templateEmail); err != nil {
		return api.Err(ctx, fiber.StatusBadRequest, "Invalid body", err)
	}

	if err := e.validate.Struct(templateEmail); err != nil {
		return api.Err(ctx, fiber.StatusBadRequest, "Invalid body", err)
	}

	template, err := e.templateService.ParseTemplate(templateName, templateEmail.Data)

	if err != nil {
		return api.Err(ctx, fiber.StatusInternalServerError, "Error parsing template", err)
	}

	mail := types.Mail{
		To:      []string{templateEmail.To},
		Subject: templateEmail.Subject,
		Sender:  config.Envs.DefaultFrom,
		Body:    template.String(),
		Type:    types.HTML,
	}

	if err != nil {
		return api.Err(ctx, fiber.StatusInternalServerError, "Error getting queue", err)
	}

	queue.Queue.Write(mail)

	api.Success(ctx, fiber.StatusCreated, "Email sent successfully")
	return nil
}
