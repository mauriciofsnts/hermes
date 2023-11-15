package smtp

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/mauriciofsnts/hermes/internal/api/controller"
	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/pauloo27/logger"
)

func SendEmail(email *types.Email) error {
	logger.Info("Sending email...")
	request, err := buildMail(*email)

	if err != nil {
		return err
	}

	var msg string

	if request.Type == types.HTML {
		msg = buildHtmlMessage(*request)
	} else {
		msg = buildTextMessage(*request)
	}

	auth := getAuth()
	addr := getAddr()

	err = smtp.SendMail(addr, auth, request.Sender, request.To, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}

func getAddr() string {
	return fmt.Sprintf("%s:%d", config.Hermes.Smtp.Host, config.Hermes.Smtp.Port)
}

func getAuth() smtp.Auth {
	return smtp.PlainAuth("", config.Hermes.Smtp.Username, config.Hermes.Smtp.Password, config.Hermes.Smtp.Host)
}

func buildMail(email types.Email) (*types.Mail, error) {
	defaultFrom := config.Hermes.DefaultFrom

	request := types.Mail{
		Sender:  defaultFrom,
		Subject: email.Subject,
		To:      []string{email.To},
		Body:    "",
	}

	if email.Body != "" {
		request.Body = email.Body
		request.Type = types.TEXT

		return &request, nil
	}

	if email.TemplateName != "" && len(email.Content) > 0 {
		controller := controller.NewTemplateController()

		buffer, err := controller.ParseTemplate(email.TemplateName, email.Content)

		if err != nil {
			return nil, err
		}

		request.Body = buffer.String()
		request.Type = types.HTML

		return &request, nil
	}

	return nil, fmt.Errorf("invalid email content")
}

func buildHtmlMessage(mail types.Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}

func buildTextMessage(mail types.Mail) string {
	msg := fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}
