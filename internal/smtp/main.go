package smtp

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/types"
	"github.com/pauloo27/logger"
)

func SendEmail(email *types.Email) error {
	logger.Info("Sending email...")

	smtpHost := config.Hermes.SmtpHost
	smtpPort := config.Hermes.SmtpPort
	smtpUsername := config.Hermes.SmtpUsername
	smtpPassword := config.Hermes.SmtpPassword
	defaultFrom := config.Hermes.DefaultFrom

	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	header := make(mail.Header)

	header["From"] = []string{defaultFrom}
	header["Subject"] = []string{email.Subject}
	header["To"] = []string{email.To}

	var msg strings.Builder

	for key, values := range header {
		msg.WriteString(key)
		msg.WriteString(": ")
		msg.WriteString(strings.Join(values, ", "))
		msg.WriteString("\r\n")
	}

	msg.WriteString("\r\n")
	msg.WriteString(email.Body)

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)

	err := smtp.SendMail(
		addr,
		auth,
		defaultFrom,
		[]string{email.To},
		[]byte(msg.String()),
	)

	if err != nil {
		return err
	}

	return nil
}
