package controller

import (
	"encoding/json"
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/mauriciofsnts/hermes/internal/api"
	"github.com/mauriciofsnts/hermes/internal/types"
)

type EmailController interface {
	SendEmail(c *fiber.Ctx) error
}

type emailControler struct {
}

func NewEmailController() *emailControler {
	return &emailControler{}
}

func (e *emailControler) SendEmail(c *fiber.Ctx) error {
	storage := c.Locals("storage").(types.Storage[types.Email])

	if storage == nil {
		return api.Err(c, fiber.StatusInternalServerError, "Failed to send email: storage is nil", nil)
	}

	email, err := e.Validation(c)

	if err != nil {
		return api.Err(c, fiber.StatusBadRequest, "Failed to send email: invalid request body", err)
	}

	err = storage.Write(*email)

	if err != nil {
		return api.Err(c, fiber.StatusInternalServerError, "Failed to send email", err)
	}

	return api.Success(c, fiber.StatusOK, "Email sent successfully")
}

func (e *emailControler) Validation(c *fiber.Ctx) (*types.Email, error) {
	payload := c.Body()

	if payload == nil {
		return nil, errors.New("invalid request body")
	}

	email := &types.Email{}

	if err := json.Unmarshal(payload, email); err != nil {
		return nil, err
	}

	if email.Body == "" && email.TemplateName == "" {
		return nil, errors.New("at least one of the fields 'body' or 'templateName' must be filled")
	}

	if email.TemplateName != "" && len(email.Content) < 1 {
		return nil, errors.New("if the field 'templateName' is filled, the field 'content' must be filled")
	}

	return email, nil
}
