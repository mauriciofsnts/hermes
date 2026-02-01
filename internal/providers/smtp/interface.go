package smtp

import (
	"github.com/mauriciofsnts/hermes/internal/types"
)

// EmailSender defines the interface for sending emails.
// This abstraction allows for easier testing and mocking of SMTP functionality.
type EmailSender interface {
	// Send sends an email with retry logic.
	// Returns an error if the email fails to send after all retry attempts.
	Send(email *types.Mail) error

	// Ping verifies the connection to the mail server.
	Ping() error

	// GetState returns the current state of the circuit breaker.
	GetState() string
}
