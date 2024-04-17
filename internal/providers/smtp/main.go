package smtp

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/types"
)

func SendEmail(email *types.Mail) error {
	var msg string

	if email.Type == types.HTML {
		msg = buildHTMLMessage(*email)
	} else {
		msg = buildTextMessage(*email)
	}

	auth := getAuth()
	addr := getAddr()

	err := smtp.SendMail(addr, auth, email.Sender, email.To, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}

func Ping() error {
	addr := getAddr()

	client, err := smtp.Dial(addr)

	if err != nil {
		return fmt.Errorf("failed to connect to smt Server: %w", err)
	}

	if err := client.Noop(); err != nil {
		return fmt.Errorf("failed to ping smt Server: %w", err)
	}

	return nil
}

func getAddr() string {
	return fmt.Sprintf("%s:%d", config.Hermes.SMTP.Host, config.Hermes.SMTP.Port)
}

func getAuth() smtp.Auth {
	return smtp.PlainAuth("", config.Hermes.SMTP.Username, config.Hermes.SMTP.Password, config.Hermes.SMTP.Host)
}

func buildHTMLMessage(mail types.Mail) string {
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
