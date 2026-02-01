package database

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DeadLetter represents a failed email notification stored for retry.
type DeadLetter struct {
	ID           uint      `gorm:"primaryKey"`
	EmailData    string    `gorm:"type:text;not null"` // JSON serialized Mail
	Error        string    `gorm:"type:text"`
	AttemptCount int       `gorm:"default:0"`
	MaxAttempts  int       `gorm:"default:5"`
	LastAttempt  time.Time `gorm:"index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Status       string `gorm:"type:varchar(20);default:'pending';index"` // pending, processing, failed, succeeded
	AppID        string `gorm:"type:varchar(100);index"`
}

// DLQService manages the Dead Letter Queue for failed emails.
type DLQService struct {
	db *gorm.DB
}

// NewDLQService creates a new Dead Letter Queue service.
// Uses SQLite by default for persistence.
func NewDLQService(dbPath string) (*DLQService, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&DeadLetter{}); err != nil {
		return nil, err
	}

	return &DLQService{db: db}, nil
}

// Store saves a failed email to the DLQ.
func (dlq *DLQService) Store(emailData string, errorMsg string, appID string) error {
	letter := &DeadLetter{
		EmailData:    emailData,
		Error:        errorMsg,
		AttemptCount: 0,
		MaxAttempts:  5,
		LastAttempt:  time.Now(),
		Status:       "pending",
		AppID:        appID,
	}

	return dlq.db.Create(letter).Error
}

// RecordAttempt updates the attempt count for a dead letter.
func (dlq *DLQService) RecordAttempt(id uint, success bool, errorMsg string) error {
	updates := map[string]interface{}{
		"attempt_count": gorm.Expr("attempt_count + 1"),
		"last_attempt":  time.Now(),
	}

	if success {
		updates["status"] = "succeeded"
	} else {
		updates["error"] = errorMsg
		// Mark as failed if max attempts reached
		var letter DeadLetter
		if err := dlq.db.First(&letter, id).Error; err != nil {
			return err
		}
		if letter.AttemptCount+1 >= letter.MaxAttempts {
			updates["status"] = "failed"
		}
	}

	return dlq.db.Model(&DeadLetter{}).Where("id = ?", id).Updates(updates).Error
}

// GetPending retrieves all pending dead letters that haven't exceeded max attempts.
func (dlq *DLQService) GetPending(limit int) ([]DeadLetter, error) {
	var letters []DeadLetter
	err := dlq.db.Where("status = ? AND attempt_count < max_attempts", "pending").
		Order("created_at ASC").
		Limit(limit).
		Find(&letters).Error
	return letters, err
}

// GetByStatus retrieves dead letters by status.
func (dlq *DLQService) GetByStatus(status string, limit int) ([]DeadLetter, error) {
	var letters []DeadLetter
	err := dlq.db.Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).
		Find(&letters).Error
	return letters, err
}

// MarkAsProcessing marks a dead letter as being processed to avoid concurrent retries.
func (dlq *DLQService) MarkAsProcessing(id uint) error {
	return dlq.db.Model(&DeadLetter{}).Where("id = ?", id).Update("status", "processing").Error
}

// ResetStuck resets dead letters stuck in processing state for more than 5 minutes.
func (dlq *DLQService) ResetStuck() error {
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
	return dlq.db.Model(&DeadLetter{}).
		Where("status = ? AND updated_at < ?", "processing", fiveMinutesAgo).
		Update("status", "pending").Error
}

// Delete removes a dead letter from the queue.
func (dlq *DLQService) Delete(id uint) error {
	return dlq.db.Delete(&DeadLetter{}, id).Error
}

// GetStats returns statistics about the DLQ.
func (dlq *DLQService) GetStats() (map[string]int64, error) {
	stats := make(map[string]int64)

	// Count by status
	var counts []struct {
		Status string
		Count  int64
	}

	err := dlq.db.Model(&DeadLetter{}).
		Select("status, count(*) as count").
		Group("status").
		Find(&counts).Error

	if err != nil {
		return nil, err
	}

	for _, c := range counts {
		stats[c.Status] = c.Count
	}

	// Total count
	var total int64
	dlq.db.Model(&DeadLetter{}).Count(&total)
	stats["total"] = total

	return stats, nil
}
