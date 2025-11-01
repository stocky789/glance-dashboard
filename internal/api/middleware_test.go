package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTokenBucketAllow(t *testing.T) {
	tb := NewTokenBucket(10, 1)

	// Should allow 10 tokens
	for i := 0; i < 10; i++ {
		if !tb.Allow(1) {
			t.Errorf("Token %d should be allowed", i+1)
		}
	}

	// 11th token should be denied
	if tb.Allow(1) {
		t.Error("11th token should be denied")
	}

	// Wait for refill
	time.Sleep(1100 * time.Millisecond)

	if !tb.Allow(1) {
		t.Error("Token should be allowed after refill")
	}
}

func TestRateLimiterGetClientIP(t *testing.T) {
	rl := NewRateLimiter()

	tests := []struct {
		name      string
		setupReq  func(*http.Request)
		expected  string
	}{
		{
			name: "RemoteAddr",
			setupReq: func(r *http.Request) {
				r.RemoteAddr = "192.168.1.1:8080"
			},
			expected: "192.168.1.1:8080",
		},
		{
			name: "X-Forwarded-For",
			setupReq: func(r *http.Request) {
				r.Header.Set("X-Forwarded-For", "10.0.0.1")
				r.RemoteAddr = "127.0.0.1:8080"
			},
			expected: "10.0.0.1",
		},
		{
			name: "X-Real-IP",
			setupReq: func(r *http.Request) {
				r.Header.Set("X-Real-IP", "172.16.0.1")
				r.RemoteAddr = "127.0.0.1:8080"
			},
			expected: "172.16.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			tt.setupReq(req)
			
			ip := rl.GetClientIP(req)
			if ip != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, ip)
			}
		})
	}
}

func TestRateLimiterAllow(t *testing.T) {
	rl := NewRateLimiter()

	for i := 0; i < 60; i++ {
		if !rl.Allow("192.168.1.1", 60) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	if rl.Allow("192.168.1.1", 60) {
		t.Error("61st request should be denied")
	}
}

func TestCORSMiddleware(t *testing.T) {
	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if origin := recorder.Header().Get("Access-Control-Allow-Origin"); origin != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin: *, got %s", origin)
	}

	if methods := recorder.Header().Get("Access-Control-Allow-Methods"); methods == "" {
		t.Error("Access-Control-Allow-Methods header not set")
	}
}

func TestCORSMiddlewareOptions(t *testing.T) {
	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	req := httptest.NewRequest("OPTIONS", "/", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS, got %d", recorder.Code)
	}
}
