package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/glance-project/glance/internal/api"
)

func TestTokenBucketAllow(t *testing.T) {
	tb := api.NewTokenBucket(10, 1)

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

func TestTokenBucketMultipleTokens(t *testing.T) {
	tb := api.NewTokenBucket(10, 1)

	if !tb.Allow(5) {
		t.Error("Should allow 5 tokens at once")
	}

	if !tb.Allow(5) {
		t.Error("Should allow another 5 tokens")
	}

	if tb.Allow(1) {
		t.Error("Should not allow 11th token")
	}
}

func TestRateLimiterGetClientIP(t *testing.T) {
	rl := api.NewRateLimiter()

	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"

	ip := rl.GetClientIP(req)
	if ip != "127.0.0.1:1234" {
		t.Errorf("Expected 127.0.0.1:1234, got %s", ip)
	}

	req.Header.Set("X-Forwarded-For", "192.168.1.1")
	ip = rl.GetClientIP(req)
	if ip != "192.168.1.1" {
		t.Errorf("Expected 192.168.1.1, got %s", ip)
	}

	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Real-IP", "10.0.0.1")
	req.RemoteAddr = "127.0.0.1:1234"

	ip = rl.GetClientIP(req)
	if ip != "10.0.0.1" {
		t.Errorf("Expected 10.0.0.1, got %s", ip)
	}
}

func TestRateLimiterAllow(t *testing.T) {
	rl := api.NewRateLimiter()

	for i := 0; i < 60; i++ {
		if !rl.Allow("192.168.1.1", 60) {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}

	if rl.Allow("192.168.1.1", 60) {
		t.Error("61st request should be denied")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	rl := api.NewRateLimiter()
	handler := api.RateLimitMiddleware(rl, 5)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "127.0.0.1:1234"
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Errorf("Request %d: expected status 200, got %d", i+1, recorder.Code)
		}
	}

	// 6th request should be rate limited
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:1234"
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status 429, got %d", recorder.Code)
	}
}

func TestCORSMiddleware(t *testing.T) {
	handler := api.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	handler := api.CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	req := httptest.NewRequest("OPTIONS", "/", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS, got %d", recorder.Code)
	}
}
