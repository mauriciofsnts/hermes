package smtp

import (
	"testing"
	"time"
)

func TestCircuitBreakerInitialization(t *testing.T) {
	cb := NewCircuitBreaker(3, 1, 30*time.Second)

	if cb.state != StateClosed {
		t.Errorf("Expected initial state to be Closed, got %d", cb.state)
	}

	if cb.GetState() != "closed" {
		t.Errorf("Expected state string to be 'closed', got %s", cb.GetState())
	}
}

func TestCircuitBreakerFailures(t *testing.T) {
	cb := NewCircuitBreaker(3, 1, 30*time.Second)

	// Registrar 2 falhas - ainda deve estar fechado
	cb.RecordFailure()
	cb.RecordFailure()
	if !cb.CanExecute() {
		t.Error("Circuit breaker should still be closed after 2 failures with threshold 3")
	}

	// 3ª falha - deve abrir
	cb.RecordFailure()
	if cb.CanExecute() {
		t.Error("Circuit breaker should be open after 3 failures")
	}

	if cb.GetState() != "open" {
		t.Errorf("Expected state to be 'open', got %s", cb.GetState())
	}
}

func TestCircuitBreakerHalfOpen(t *testing.T) {
	cb := NewCircuitBreaker(3, 1, 100*time.Millisecond)

	// Abrir o circuit breaker
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}

	// Não pode executar quando está aberto
	if cb.CanExecute() {
		t.Error("Circuit breaker should not allow execution in Open state")
	}

	// Aguardar timeout
	time.Sleep(150 * time.Millisecond)

	// Deve estar em Half-Open
	if cb.CanExecute() != true {
		t.Error("Circuit breaker should allow execution after timeout (Half-Open state)")
	}

	if cb.GetState() != "half-open" {
		t.Errorf("Expected state to be 'half-open', got %s", cb.GetState())
	}
}

func TestCircuitBreakerRecovery(t *testing.T) {
	cb := NewCircuitBreaker(3, 1, 100*time.Millisecond)

	// Abrir o circuit breaker
	for i := 0; i < 3; i++ {
		cb.RecordFailure()
	}

	// Aguardar e ir para Half-Open
	time.Sleep(150 * time.Millisecond)
	cb.CanExecute()

	// Registrar sucesso
	cb.RecordSuccess()

	// Deve estar fechado novamente
	if cb.GetState() != "closed" {
		t.Errorf("Expected state to be 'closed' after recovery, got %s", cb.GetState())
	}

	if cb.failureCount != 0 {
		t.Errorf("Expected failure count to be reset, got %d", cb.failureCount)
	}
}

func TestRetryConfig(t *testing.T) {
	config := DefaultRetryConfig()

	if config.MaxAttempts != 3 {
		t.Errorf("Expected MaxAttempts to be 3, got %d", config.MaxAttempts)
	}

	if config.InitialDelay != 1*time.Second {
		t.Errorf("Expected InitialDelay to be 1s, got %v", config.InitialDelay)
	}

	if config.BackoffMultiplier != 2.0 {
		t.Errorf("Expected BackoffMultiplier to be 2.0, got %f", config.BackoffMultiplier)
	}
}

func TestBackoffDurationCalculation(t *testing.T) {
	tests := []struct {
		attempt    int
		initial    time.Duration
		multiplier float64
		max        time.Duration
		expected   time.Duration
	}{
		{0, 1 * time.Second, 2.0, 30 * time.Second, 1 * time.Second},
		{1, 1 * time.Second, 2.0, 30 * time.Second, 2 * time.Second},
		{2, 1 * time.Second, 2.0, 30 * time.Second, 4 * time.Second},
		{3, 1 * time.Second, 2.0, 30 * time.Second, 8 * time.Second},
		{10, 1 * time.Second, 2.0, 30 * time.Second, 30 * time.Second}, // capped at max
	}

	for _, tt := range tests {
		result := CalculateBackoffDuration(tt.attempt, tt.initial, tt.multiplier, tt.max)
		if result != tt.expected {
			t.Errorf("Attempt %d: expected %v, got %v", tt.attempt, tt.expected, result)
		}
	}
}
