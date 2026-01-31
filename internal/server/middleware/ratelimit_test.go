package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiterAllows(t *testing.T) {
	limiter := NewRateLimiter(10) // 10 requests per hour
	ip := "192.168.1.1"

	// Deve permitir 10 requisições
	for i := 0; i < 10; i++ {
		if !limiter.IsAllowed(ip) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	// 11ª deve ser bloqueada
	if limiter.IsAllowed(ip) {
		t.Error("Request 11 should be blocked")
	}
}

func TestRateLimiterPerIP(t *testing.T) {
	limiter := NewRateLimiter(2) // 2 requests per hour
	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	// IP1: 2 requisições
	if !limiter.IsAllowed(ip1) {
		t.Error("IP1 request 1 should be allowed")
	}
	if !limiter.IsAllowed(ip1) {
		t.Error("IP1 request 2 should be allowed")
	}
	if limiter.IsAllowed(ip1) {
		t.Error("IP1 request 3 should be blocked")
	}

	// IP2: deve ter seu próprio limite
	if !limiter.IsAllowed(ip2) {
		t.Error("IP2 request 1 should be allowed")
	}
	if !limiter.IsAllowed(ip2) {
		t.Error("IP2 request 2 should be allowed")
	}
	if limiter.IsAllowed(ip2) {
		t.Error("IP2 request 3 should be blocked")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	limiter := NewRateLimiter(2)
	handler := RateLimitMiddleware(limiter)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	// Primeira requisição
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "192.168.1.1:5000"
	w1 := httptest.NewRecorder()
	handler.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Errorf("First request should be OK, got %d", w1.Code)
	}

	// Segunda requisição
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "192.168.1.1:5000"
	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Errorf("Second request should be OK, got %d", w2.Code)
	}

	// Terceira requisição deve ser bloqueada
	req3 := httptest.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "192.168.1.1:5000"
	w3 := httptest.NewRecorder()
	handler.ServeHTTP(w3, req3)

	if w3.Code != http.StatusTooManyRequests {
		t.Errorf("Third request should be TooManyRequests, got %d", w3.Code)
	}
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*http.Request)
		expected string
	}{
		{
			name: "X-Forwarded-For header",
			setup: func(r *http.Request) {
				r.Header.Set("X-Forwarded-For", "10.0.0.1")
			},
			expected: "10.0.0.1",
		},
		{
			name: "X-Real-IP header",
			setup: func(r *http.Request) {
				r.Header.Set("X-Real-IP", "10.0.0.2")
			},
			expected: "10.0.0.2",
		},
		{
			name: "RemoteAddr",
			setup: func(r *http.Request) {
				r.RemoteAddr = "10.0.0.3:5000"
			},
			expected: "10.0.0.3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			tt.setup(req)
			ip := getClientIP(req)

			if ip != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, ip)
			}
		})
	}
}

func TestRateLimiterCleanup(t *testing.T) {
	// Usar uma janela muito curta para teste
	limiter := &RateLimiter{
		requestCounts: make(map[string][]time.Time),
		limit:         10,
		window:        100 * time.Millisecond, // Janela de 100ms
	}

	ip := "192.168.1.1"
	limiter.IsAllowed(ip)

	// Verificar se foi armazenado
	limiter.mu.RLock()
	initialCount := len(limiter.requestCounts)
	limiter.mu.RUnlock()

	if initialCount != 1 {
		t.Errorf("Expected 1 IP in map, got %d", initialCount)
	}

	// Aguardar janela expirar
	time.Sleep(150 * time.Millisecond)

	// Executar cleanup
	limiter.cleanup()

	// Verificar se foi removido
	limiter.mu.RLock()
	finalCount := len(limiter.requestCounts)
	limiter.mu.RUnlock()

	if finalCount != 0 {
		t.Errorf("Expected 0 IPs after cleanup, got %d", finalCount)
	}
}
