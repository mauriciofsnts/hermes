package http

import (
	"fmt"

	"net/smtp"

	"github.com/gofiber/fiber/v2"

	"github.com/mauriciofsnts/hermes/internal/config"
)

type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func Listen() error {
	app := fiber.New()

	smtpHost := config.Hermes.SmtpHost
	smtpPort := config.Hermes.SmtpPort
	smtpUsername := config.Hermes.SmtpUsername
	smtpPassword := config.Hermes.SmtpPassword
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	defaultFrom := config.Hermes.DefaultFrom

	app.Post("/api/send-email", func(c *fiber.Ctx) error {
		var email Email
		if err := c.BodyParser(&email); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"Error": "invalid request body",
			})
		}

		auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

		// Crie a mensagem de e-mail
		msg := []byte(
			"From: " + defaultFrom + "\r\n" +
				"To: " + email.To + "\r\n" +
				"Subject: " + email.Subject + "\r\n" +
				"\r\n" +
				email.Body + "\r\n",
		)

		err := smtp.SendMail(
			addr,
			auth,
			defaultFrom,
			[]string{email.To},
			msg,
		)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"Error": "failed to send email: " + err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"Message": "email sent successfully",
		})
	})

	app.Listen(":8080")

	return nil
}
