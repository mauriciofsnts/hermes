package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/mauriciofsnts/hermes/internal/metrics"
	"github.com/mauriciofsnts/hermes/internal/server/api"
	"github.com/redis/go-redis/v9"
)

// RedisRateLimiter implements distributed rate limiting using Redis.
type RedisRateLimiter struct {
	redis      *redis.Client
	limitPerIP int
	window     time.Duration
}

// NewRedisRateLimiter creates a new Redis-backed rate limiter.
func NewRedisRateLimiter(redisClient *redis.Client, limitPerIP int, window time.Duration) *RedisRateLimiter {
	return &RedisRateLimiter{
		redis:      redisClient,
		limitPerIP: limitPerIP,
		window:     window,
	}
}

// Allow checks if a request from the given IP is allowed.
func (rl *RedisRateLimiter) Allow(ip string) bool {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limit:ip:%s", ip)

	// Use Redis INCR with EXPIRE for atomic rate limiting
	count, err := rl.redis.Incr(ctx, key).Result()
	if err != nil {
		// On error, allow the request (fail open)
		return true
	}

	// Set expiration on first request
	if count == 1 {
		rl.redis.Expire(ctx, key, rl.window)
	}

	return count <= int64(rl.limitPerIP)
}

// GetRemaining returns the number of remaining requests for an IP.
func (rl *RedisRateLimiter) GetRemaining(ip string) int {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limit:ip:%s", ip)

	count, err := rl.redis.Get(ctx, key).Int()
	if err != nil {
		return rl.limitPerIP
	}

	remaining := rl.limitPerIP - count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetTTL returns the time until the rate limit resets for an IP.
func (rl *RedisRateLimiter) GetTTL(ip string) time.Duration {
	ctx := context.Background()
	key := fmt.Sprintf("rate_limit:ip:%s", ip)

	ttl, err := rl.redis.TTL(ctx, key).Result()
	if err != nil || ttl < 0 {
		return rl.window
	}
	return ttl
}

// RedisRateLimitMiddleware creates a middleware for distributed rate limiting.
func RedisRateLimitMiddleware(rl *RedisRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)

			if !rl.Allow(ip) {
				// Record rate limit violation in metrics
				metrics.APIKeyRateLimit.WithLabelValues(ip, "ip").Inc()

				ttl := rl.GetTTL(ip)
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.limitPerIP))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(ttl).Unix(), 10))
				w.Header().Set("Retry-After", strconv.Itoa(int(ttl.Seconds())))
				w.WriteHeader(http.StatusTooManyRequests)

				response := api.TooManyRequestsErr("Rate limit exceeded")
				encoder := json.NewEncoder(w)
				_ = encoder.Encode(response.Body)
				return
			}

			// Set rate limit headers
			remaining := rl.GetRemaining(ip)
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(rl.limitPerIP))
			w.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(rl.window).Unix(), 10))

			next.ServeHTTP(w, r)
		})
	}
}
