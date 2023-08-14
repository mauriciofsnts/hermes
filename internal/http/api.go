package http

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mauriciofsnts/hermes/internal/config"
)

type Email struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func SendEmail(c *fiber.Ctx) error {
	smtpHost := config.Hermes.SmtpHost
	smtpPort := config.Hermes.SmtpPort
	smtpUsername := config.Hermes.SmtpUsername
	smtpPassword := config.Hermes.SmtpPassword

	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	defaultFrom := config.Hermes.DefaultFrom
	allowedOrigin := config.Hermes.AllowedOrigin

	var email Email

	if allowedOrigin != "" {
		origin := c.Get("Origin")

		if origin != allowedOrigin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Origin not allowed",
			})
		}
	}

	if err := c.BodyParser(&email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	header := make(mail.Header)

	header["From"] = []string{defaultFrom}
	header["To"] = []string{email.To}
	header["Subject"] = []string{email.Subject}

	var msg strings.Builder

	for key, values := range header {
		msg.WriteString(key)
		msg.WriteString(": ")
		msg.WriteString(strings.Join(values, ", "))
		msg.WriteString("\r\n")
	}

	msg.WriteString("\r\n")
	msg.WriteString(email.Body)

	err := smtp.SendMail(
		addr,
		auth,
		defaultFrom,
		[]string{email.To},
		[]byte(msg.String()),
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send email: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Email sent successfully",
	})
}

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hermes is up and running",
	})
}
