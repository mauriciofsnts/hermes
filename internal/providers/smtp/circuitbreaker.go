package smtp

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type CircuitBreakerState int

const (
	StateClosed CircuitBreakerState = iota
	StateOpen
	StateHalfOpen
)

type CircuitBreaker struct {
	mu               sync.RWMutex
	state            CircuitBreakerState
	failureCount     int
	failureThreshold int
	successCount     int
	successThreshold int
	lastFailureTime  time.Time
	timeout          time.Duration
}

func NewCircuitBreaker(failureThreshold int, successThreshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:            StateClosed,
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		timeout:          timeout,
	}
}

func (cb *CircuitBreaker) CanExecute() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = StateHalfOpen
			cb.successCount = 0
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		cb.failureCount = 0
	case StateHalfOpen:
		cb.successCount++
		if cb.successCount >= cb.successThreshold {
			cb.state = StateClosed
			cb.failureCount = 0
		}
	}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.lastFailureTime = time.Now()
	switch cb.state {
	case StateClosed:
		cb.failureCount++
		if cb.failureCount >= cb.failureThreshold {
			cb.state = StateOpen
		}
	case StateHalfOpen:
		cb.state = StateOpen
	}
}

func (cb *CircuitBreaker) GetState() string {
	cb.mu.RLock()
	defer cb.mu.RUnlock()

	switch cb.state {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	}
	return "unknown"
}

// RetryConfig configures retry behavior with exponential backoff.
type RetryConfig struct {
	MaxAttempts       int
	InitialDelay      time.Duration
	MaxDelay          time.Duration
	BackoffMultiplier float64
}

// DefaultRetryConfig returns the default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:       3,
		InitialDelay:      1 * time.Second,
		MaxDelay:          30 * time.Second,
		BackoffMultiplier: 2.0,
	}
}

// ExecuteWithRetry executes a function with retry and exponential backoff.
func ExecuteWithRetry(fn func() error, config RetryConfig) error {
	var lastErr error
	delay := config.InitialDelay

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		if attempt == config.MaxAttempts-1 {
			break
		}

		time.Sleep(delay)

		nextDelay := time.Duration(float64(delay) * config.BackoffMultiplier)
		if nextDelay > config.MaxDelay {
			nextDelay = config.MaxDelay
		}
		delay = nextDelay
	}

	return fmt.Errorf("after %d attempts: %w", config.MaxAttempts, lastErr)
}

// CalculateBackoffDuration calculates the delay based on the attempt number.
func CalculateBackoffDuration(attempt int, initialDelay time.Duration, multiplier float64, maxDelay time.Duration) time.Duration {
	delay := time.Duration(float64(initialDelay) * math.Pow(multiplier, float64(attempt)))
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}
