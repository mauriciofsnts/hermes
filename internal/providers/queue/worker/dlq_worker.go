package worker

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/mauriciofsnts/hermes/internal/providers/database"
	"github.com/mauriciofsnts/hermes/internal/providers/smtp"
	"github.com/mauriciofsnts/hermes/internal/types"
)

// DLQWorker processes failed emails from the Dead Letter Queue.
type DLQWorker struct {
	dlq      *database.DLQService
	smtp     smtp.EmailSender
	interval time.Duration
}

// NewDLQWorker creates a new DLQ worker.
func NewDLQWorker(dlq *database.DLQService, smtpProvider smtp.EmailSender, interval time.Duration) *DLQWorker {
	return &DLQWorker{
		dlq:      dlq,
		smtp:     smtpProvider,
		interval: interval,
	}
}

// Start begins processing the DLQ periodically.
func (w *DLQWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	slog.Info("DLQ Worker started", "interval", w.interval)

	for {
		select {
		case <-ctx.Done():
			slog.Info("DLQ Worker stopped")
			return
		case <-ticker.C:
			w.processBatch(ctx)
		}
	}
}

// processBatch processes a batch of pending dead letters.
func (w *DLQWorker) processBatch(ctx context.Context) {
	// Reset stuck items first
	if err := w.dlq.ResetStuck(); err != nil {
		slog.Error("Failed to reset stuck DLQ items", "error", err)
	}

	// Get pending items
	letters, err := w.dlq.GetPending(10)
	if err != nil {
		slog.Error("Failed to get pending DLQ items", "error", err)
		return
	}

	if len(letters) == 0 {
		return
	}

	slog.Info("Processing DLQ batch", "count", len(letters))

	for _, letter := range letters {
		select {
		case <-ctx.Done():
			return
		default:
			w.processItem(letter)
		}
	}
}

// processItem attempts to resend a single dead letter.
func (w *DLQWorker) processItem(letter database.DeadLetter) {
	// Mark as processing
	if err := w.dlq.MarkAsProcessing(letter.ID); err != nil {
		slog.Error("Failed to mark DLQ item as processing", "id", letter.ID, "error", err)
		return
	}

	// Deserialize email
	var email types.Mail
	if err := json.Unmarshal([]byte(letter.EmailData), &email); err != nil {
		slog.Error("Failed to deserialize email from DLQ", "id", letter.ID, "error", err)
		w.dlq.RecordAttempt(letter.ID, false, "deserialization error: "+err.Error())
		return
	}

	// Attempt to send
	slog.Info("Retrying email from DLQ",
		"id", letter.ID,
		"attempt", letter.AttemptCount+1,
		"to", email.To,
	)

	err := w.smtp.Send(&email)
	if err != nil {
		slog.Warn("DLQ retry failed",
			"id", letter.ID,
			"attempt", letter.AttemptCount+1,
			"error", err,
		)
		w.dlq.RecordAttempt(letter.ID, false, err.Error())
	} else {
		slog.Info("DLQ retry succeeded", "id", letter.ID)
		w.dlq.RecordAttempt(letter.ID, true, "")
	}
}
