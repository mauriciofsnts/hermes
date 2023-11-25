package controller

import (
	"errors"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/mauriciofsnts/hermes/internal/api"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type EmailControllerInterface interface {
	SendPlainTextEmail(c *fiber.Ctx) error
	SendTemplateEmail(c *fiber.Ctx) error
}

type EmailControler struct {
	EmailControllerInterface
	validate *validator.Validate
}

func NewEmailController() *EmailControler {
	return &EmailControler{
		validate: validator.New(),
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
		Sender:  config.Hermes.DefaultFrom,
		Body:    bodyEmail.Body,
		Type:    types.TEXT,
	}

	queue, err := e.getQueue(c)

	if err != nil {
		return api.Err(c, fiber.StatusInternalServerError, "Error getting queue", err)
	}

	queue.Write(mail)

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

	templateController := NewTemplateController()

	template, err := templateController.ParseTemplate(templateName, templateEmail.Data)

	if err != nil {
		return api.Err(ctx, fiber.StatusInternalServerError, "Error parsing template", err)
	}

	slog.Any("template", template.String())

	mail := types.Mail{
		To:      []string{templateEmail.To},
		Subject: templateEmail.Subject,
		Sender:  config.Hermes.DefaultFrom,
		Body:    template.String(),
		Type:    types.HTML,
	}

	queue, err := e.getQueue(ctx)

	if err != nil {
		return api.Err(ctx, fiber.StatusInternalServerError, "Error getting queue", err)
	}

	queue.Write(mail)

	api.Success(ctx, fiber.StatusCreated, "Email sent successfully")
	return nil
}

func (e *EmailControler) getQueue(c *fiber.Ctx) (types.Queue[types.Mail], error) {
	queue := c.Locals("queue").(types.Queue[types.Mail])

	if queue == nil {
		return nil, errors.New("queue not found")
	}

	return queue, nil
}
