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

// SMTPProvider implements the EmailSender interface for SMTP email delivery.
type SMTPProvider struct {
	breaker *CircuitBreaker
}

// NewSMTPProvider creates a new SMTP provider with a circuit breaker.
// The circuit breaker opens after 3 failures, half-opens after 1 minute.
func NewSMTPProvider() *SMTPProvider {
	return &SMTPProvider{
		breaker: NewCircuitBreaker(3, 1, 30*time.Second),
	}
}

// Send sends an email with automatic retry logic.
// It respects the circuit breaker state to prevent cascading failures.
func (sp *SMTPProvider) Send(email *types.Mail) error {
	if !sp.breaker.CanExecute() {
		return fmt.Errorf("SMTP circuit breaker is %s", sp.breaker.GetState())
	}

	retryConfig := DefaultRetryConfig()
	err := ExecuteWithRetry(func() error {
		return sp.sendWithoutRetry(email)
	}, retryConfig)

	if err != nil {
		sp.breaker.RecordFailure()
		slog.Error("Failed to send email after retries", "to", email.To, "error", err)
		return err
	}

	sp.breaker.RecordSuccess()
	return nil
}

// sendWithoutRetry performs a single SMTP send attempt without retries.
func (sp *SMTPProvider) sendWithoutRetry(email *types.Mail) error {
	msg := buildHTMLMessage(*email)

	auth := getAuth()
	addr := getAddr()

	return smtp.SendMail(addr, auth, email.Sender, email.To, []byte(msg))
}

// Ping verifies the connection to the SMTP server.
func (sp *SMTPProvider) Ping() error {
	addr := getAddr()

	client, err := smtp.Dial(addr)

	if err != nil {
		slog.Error("failed to connect to smtp server", "error", err)
		return fmt.Errorf("failed to connect to smtp server: %w", err)
	}

	if err := client.Noop(); err != nil {
		slog.Error("failed to ping smtp server", "error", err)
		return fmt.Errorf("failed to ping smtp server: %w", err)
	}

	return nil
}

// GetState returns the current state of the circuit breaker.
func (sp *SMTPProvider) GetState() string {
	return sp.breaker.GetState()
}

// Legacy package-level functions for backward compatibility
var defaultProvider *SMTPProvider

func init() {
	defaultProvider = NewSMTPProvider()
}

// SendEmail sends an email using the default SMTP provider.
// Deprecated: Use NewSMTPProvider().Send() instead.
func SendEmail(email *types.Mail) error {
	return defaultProvider.Send(email)
}

// Ping verifies the connection to the SMTP server using the default provider.
// Deprecated: Use NewSMTPProvider().Ping() instead.
func Ping() error {
	return defaultProvider.Ping()
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
