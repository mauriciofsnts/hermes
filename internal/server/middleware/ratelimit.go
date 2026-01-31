package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/mauriciofsnts/hermes/internal/server/api"
)

type RateLimiter struct {
	requestCounts map[string][]time.Time
	limit         int
	window        time.Duration
	mu            sync.RWMutex
}

func NewRateLimiter(requestsPerHour int) *RateLimiter {
	limiter := &RateLimiter{
		requestCounts: make(map[string][]time.Time),
		limit:         requestsPerHour,
		window:        time.Hour,
	}

	// Limpar requisições antigas periodicamente
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			limiter.cleanup()
		}
	}()

	return limiter
}

// IsAllowed verifica se a requisição é permitida para o IP
func (rl *RateLimiter) IsAllowed(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	requests := rl.requestCounts[ip]

	// Remover requisições fora da janela de tempo
	var validRequests []time.Time
	for _, req := range requests {
		if now.Sub(req) < rl.window {
			validRequests = append(validRequests, req)
		}
	}

	// Verificar se limite foi atingido
	if len(validRequests) >= rl.limit {
		return false
	}

	// Adicionar nova requisição
	validRequests = append(validRequests, now)
	rl.requestCounts[ip] = validRequests
	return true
}

// cleanup remove entradas antigas de IPs
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, requests := range rl.requestCounts {
		var validRequests []time.Time
		for _, req := range requests {
			if now.Sub(req) < rl.window {
				validRequests = append(validRequests, req)
			}
		}

		if len(validRequests) == 0 {
			delete(rl.requestCounts, ip)
		} else {
			rl.requestCounts[ip] = validRequests
		}
	}
}

// RateLimitMiddleware retorna um middleware HTTP que aplica rate limiting
func RateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)

			if !limiter.IsAllowed(ip) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				resp := api.TooManyRequestsErr("rate limit exceeded")
				encoder := json.NewEncoder(w)
				_ = encoder.Encode(resp.Body)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// getClientIP extrai o IP do cliente da requisição
func getClientIP(r *http.Request) string {
	// Verificar X-Forwarded-For (quando atrás de proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip, _, err := net.SplitHostPort(xff); err == nil {
			return ip
		}
		return xff
	}

	// Verificar X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
