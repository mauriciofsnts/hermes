package notify

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/ctx"
	"github.com/mauriciofsnts/hermes/internal/providers/queue"
	"github.com/mauriciofsnts/hermes/internal/providers/template"
	"github.com/mauriciofsnts/hermes/internal/server/helper"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type EmailControllerInterface interface {
	SendPlainTextEmail(c *fiber.Ctx) error
	SendTemplateEmail(c *fiber.Ctx) error
}

type EmailControler struct {
	EmailControllerInterface
	validate         *validator.Validate
	templateProvider template.TemplateProvider
}

func NewEmailController() *EmailControler {
	return &EmailControler{
		validate:         validator.New(),
		templateProvider: template.NewTemplateService(),
	}
}

func (e *EmailControler) SendPlainTextEmail(c *fiber.Ctx) error {
	var bodyEmail types.PlainTextEmail
	ctxProviders := c.Locals("providers").(*ctx.Providers)

	if ctxProviders == nil {
		return helper.Err(c, fiber.StatusInternalServerError, "Providers not found", nil)
	}

	queue := ctxProviders.Queue

	if err := c.BodyParser(&bodyEmail); err != nil {
		return helper.Err(c, fiber.StatusBadRequest, "Invalid body", err)
	}

	if err := e.validate.Struct(bodyEmail); err != nil {
		return helper.Err(c, fiber.StatusBadRequest, "Invalid body", err)
	}

	mail := types.Mail{
		To:      []string{bodyEmail.To},
		Subject: bodyEmail.Subject,
		Sender:  config.Hermes.SMTP.Sender,
		Body:    bodyEmail.Body,
		Type:    types.TEXT,
	}

	queue.Write(mail)

	helper.Success(c, fiber.StatusCreated, "Email sent successfully")
	return nil
}

func (e *EmailControler) SendTemplateEmail(ctx *fiber.Ctx) error {
	var templateEmail types.TemplateEmail

	templateName := ctx.Params("slug")

	if templateName == "" {
		return helper.Err(ctx, fiber.StatusBadRequest, "Invalid template name", nil)
	}

	if err := ctx.BodyParser(&templateEmail); err != nil {
		return helper.Err(ctx, fiber.StatusBadRequest, "Invalid body", err)
	}

	if err := e.validate.Struct(templateEmail); err != nil {
		return helper.Err(ctx, fiber.StatusBadRequest, "Invalid body", err)
	}

	template, err := e.templateProvider.ParseTemplate(templateName, templateEmail.Data)

	if err != nil {
		return helper.Err(ctx, fiber.StatusInternalServerError, "Error parsing template", err)
	}

	mail := types.Mail{
		To:      []string{templateEmail.To},
		Subject: templateEmail.Subject,
		Sender:  config.Hermes.SMTP.Sender,
		Body:    template.String(),
		Type:    types.HTML,
	}

	queue.Queue.Write(mail)

	helper.Success(ctx, fiber.StatusCreated, "Email sent successfully")
	return nil
}
