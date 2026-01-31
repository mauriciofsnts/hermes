package smtp

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"
	"time"

	"github.com/mauriciofsnts/hermes/internal/config"
	"github.com/mauriciofsnts/hermes/internal/types"
)

var smtpCircuitBreaker = NewCircuitBreaker(3, 1, 30*time.Second)

func SendEmail(email *types.Mail) error {
	if !smtpCircuitBreaker.CanExecute() {
		return fmt.Errorf("SMTP circuit breaker is %s", smtpCircuitBreaker.GetState())
	}

	retryConfig := DefaultRetryConfig()
	err := ExecuteWithRetry(func() error {
		return sendEmailWithoutRetry(email)
	}, retryConfig)

	if err != nil {
		smtpCircuitBreaker.RecordFailure()
		slog.Error("Failed to send email after retries", "to", email.To, "error", err)
		return err
	}

	smtpCircuitBreaker.RecordSuccess()
	return nil
}

func sendEmailWithoutRetry(email *types.Mail) error {
	msg := buildHTMLMessage(*email)

	auth := getAuth()
	addr := getAddr()

	return smtp.SendMail(addr, auth, email.Sender, email.To, []byte(msg))
}

func Ping() error {
	addr := getAddr()

	client, err := smtp.Dial(addr)

	if err != nil {
		slog.Error("failed to connect to smt Server", "error", err)
		return fmt.Errorf("failed to connect to smt Server: %w", err)
	}

	if err := client.Noop(); err != nil {
		slog.Error("failed to ping smt Server", "error", err)
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
