package smtp

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCircuitBreaker implements a distributed circuit breaker using Redis.
// This allows multiple instances of the application to share circuit breaker state.
type RedisCircuitBreaker struct {
	redis               *redis.Client
	keyPrefix           string
	failureThreshold    int
	successThreshold    int
	timeout             time.Duration
	halfOpenMaxAttempts int
}

// NewRedisCircuitBreaker creates a new Redis-backed circuit breaker.
func NewRedisCircuitBreaker(redisClient *redis.Client, keyPrefix string, failureThreshold int, successThreshold int, timeout time.Duration) *RedisCircuitBreaker {
	return &RedisCircuitBreaker{
		redis:               redisClient,
		keyPrefix:           keyPrefix,
		failureThreshold:    failureThreshold,
		successThreshold:    successThreshold,
		timeout:             timeout,
		halfOpenMaxAttempts: 1,
	}
}

// CanExecute checks if the circuit breaker allows execution.
func (rcb *RedisCircuitBreaker) CanExecute() bool {
	ctx := context.Background()
	state := rcb.getState(ctx)

	switch state {
	case "closed":
		return true
	case "open":
		// Check if timeout has elapsed
		lastFailureTime, err := rcb.redis.Get(ctx, rcb.key("last_failure")).Result()
		if err != nil {
			return true // If we can't read, allow execution
		}

		lastFailure, err := time.Parse(time.RFC3339, lastFailureTime)
		if err != nil {
			return true
		}

		if time.Since(lastFailure) > rcb.timeout {
			// Transition to half-open
			rcb.setState(ctx, "half-open")
			rcb.redis.Set(ctx, rcb.key("half_open_attempts"), "0", 0)
			return true
		}
		return false
	case "half-open":
		// Check if we've exceeded half-open attempts
		attempts, _ := rcb.redis.Get(ctx, rcb.key("half_open_attempts")).Int()
		return attempts < rcb.halfOpenMaxAttempts
	default:
		return true
	}
}

// RecordSuccess records a successful execution.
func (rcb *RedisCircuitBreaker) RecordSuccess() {
	ctx := context.Background()
	state := rcb.getState(ctx)

	switch state {
	case "half-open":
		// Increment success count
		successes := rcb.redis.Incr(ctx, rcb.key("successes")).Val()
		if int(successes) >= rcb.successThreshold {
			// Transition back to closed
			rcb.setState(ctx, "closed")
			rcb.redis.Del(ctx, rcb.key("failures"), rcb.key("successes"), rcb.key("half_open_attempts"))
		}
	case "closed":
		// Reset failure count on success
		rcb.redis.Set(ctx, rcb.key("failures"), "0", 0)
	}
}

// RecordFailure records a failed execution.
func (rcb *RedisCircuitBreaker) RecordFailure() {
	ctx := context.Background()
	state := rcb.getState(ctx)

	rcb.redis.Set(ctx, rcb.key("last_failure"), time.Now().Format(time.RFC3339), rcb.timeout*2)

	switch state {
	case "closed":
		failures := rcb.redis.Incr(ctx, rcb.key("failures")).Val()
		if int(failures) >= rcb.failureThreshold {
			rcb.setState(ctx, "open")
		}
	case "half-open":
		// Transition back to open
		rcb.setState(ctx, "open")
		rcb.redis.Del(ctx, rcb.key("successes"), rcb.key("half_open_attempts"))
	}

	// Increment half-open attempts if in half-open state
	if state == "half-open" {
		rcb.redis.Incr(ctx, rcb.key("half_open_attempts"))
	}
}

// GetState returns the current state of the circuit breaker.
func (rcb *RedisCircuitBreaker) GetState() string {
	return rcb.getState(context.Background())
}

// getState retrieves the current state from Redis.
func (rcb *RedisCircuitBreaker) getState(ctx context.Context) string {
	state, err := rcb.redis.Get(ctx, rcb.key("state")).Result()
	if err != nil {
		return "closed" // Default to closed if not set
	}
	return state
}

// setState sets the current state in Redis.
func (rcb *RedisCircuitBreaker) setState(ctx context.Context, state string) {
	rcb.redis.Set(ctx, rcb.key("state"), state, 0)
}

// key generates a Redis key with the configured prefix.
func (rcb *RedisCircuitBreaker) key(suffix string) string {
	return fmt.Sprintf("%s:%s", rcb.keyPrefix, suffix)
}

// GetStats returns statistics about the circuit breaker.
func (rcb *RedisCircuitBreaker) GetStats() map[string]string {
	ctx := context.Background()
	stats := make(map[string]string)

	stats["state"] = rcb.getState(ctx)

	failures, _ := rcb.redis.Get(ctx, rcb.key("failures")).Result()
	stats["failures"] = failures

	successes, _ := rcb.redis.Get(ctx, rcb.key("successes")).Result()
	stats["successes"] = successes

	lastFailure, _ := rcb.redis.Get(ctx, rcb.key("last_failure")).Result()
	stats["last_failure"] = lastFailure

	return stats
}

// Reset resets the circuit breaker to its initial state.
func (rcb *RedisCircuitBreaker) Reset() {
	ctx := context.Background()
	rcb.redis.Del(ctx,
		rcb.key("state"),
		rcb.key("failures"),
		rcb.key("successes"),
		rcb.key("half_open_attempts"),
		rcb.key("last_failure"),
	)
}

// Sync ensures the circuit breaker state is synchronized (no-op for Redis implementation).
func (rcb *RedisCircuitBreaker) Sync() error {
	// Redis handles synchronization automatically
	return nil
}
