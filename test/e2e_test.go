package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"
)

type APIClient struct {
	baseURL string
	client  *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *APIClient) Get(path string) (*http.Response, error) {
	return c.client.Get(c.baseURL + path)
}

func (c *APIClient) Post(path string, data interface{}) (*http.Response, error) {
	body, _ := json.Marshal(data)
	return c.client.Post(c.baseURL+path, "application/json", bytes.NewBuffer(body))
}

func TestAPIEndpoints(t *testing.T) {
	client := NewAPIClient("http://localhost:8080")

	tests := []struct {
		name   string
		method string
		path   string
	}{
		{"Health Check", "GET", "/api/v1/health"},
		{"Metrics", "GET", "/api/v1/metrics"},
		{"WebSocket", "GET", "/api/ws"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := client.Get(test.path)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				t.Errorf("Expected 2xx or 3xx, got %d", resp.StatusCode)
			}
		})
	}
}

func TestLoadTest(t *testing.T) {
	client := NewAPIClient("http://localhost:8080")
	numRequests := 100
	concurrency := 10

	results := make([]time.Duration, numRequests)
	var mu sync.Mutex
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrency)

	start := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(index int) {
			defer wg.Done()
			defer func() { <-semaphore }()

			requestStart := time.Now()
			resp, err := client.Get("/api/v1/metrics")
			duration := time.Since(requestStart)

			if err != nil {
				t.Errorf("Request %d failed: %v", index, err)
				return
			}
			defer resp.Body.Close()

			mu.Lock()
			results[index] = duration
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	totalDuration := time.Since(start)

	// Calculate statistics
	var totalTime time.Duration
	minTime := time.Duration(1<<63 - 1)
	maxTime := time.Duration(0)

	for _, d := range results {
		totalTime += d
		if d < minTime {
			minTime = d
		}
		if d > maxTime {
			maxTime = d
		}
	}

	avgTime := totalTime / time.Duration(numRequests)
	throughput := float64(numRequests) / totalDuration.Seconds()

	t.Logf("Load Test Results:")
	t.Logf("  Total Requests: %d", numRequests)
	t.Logf("  Total Duration: %v", totalDuration)
	t.Logf("  Avg Time: %v", avgTime)
	t.Logf("  Min Time: %v", minTime)
	t.Logf("  Max Time: %v", maxTime)
	t.Logf("  Throughput: %.2f req/s", throughput)
}

func TestConcurrentRequests(t *testing.T) {
	client := NewAPIClient("http://localhost:8080")
	numGoroutines := 50
	requestsPerGoroutine := 20

	var wg sync.WaitGroup
	errorCount := 0
	var mu sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < requestsPerGoroutine; j++ {
				resp, err := client.Get("/api/v1/metrics")
				if err != nil {
					mu.Lock()
					errorCount++
					mu.Unlock()
					t.Logf("Goroutine %d request %d failed: %v", id, j, err)
					continue
				}
				resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()

	if errorCount > 0 {
		t.Errorf("Expected 0 errors, got %d", errorCount)
	}
}

func TestRateLimiting(t *testing.T) {
	client := NewAPIClient("http://localhost:8080")
	maxRequests := 10
	var successCount int
	var rateLimitedCount int

	for i := 0; i < maxRequests+5; i++ {
		resp, err := client.Get("/api/v1/metrics")
		if err != nil {
			t.Fatalf("Request failed: %v", err)
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			rateLimitedCount++
		} else if resp.StatusCode == http.StatusOK {
			successCount++
		}
		resp.Body.Close()

		time.Sleep(100 * time.Millisecond)
	}

	t.Logf("Rate Limiting Test:")
	t.Logf("  Successful Requests: %d", successCount)
	t.Logf("  Rate Limited Requests: %d", rateLimitedCount)
}

func TestErrorHandling(t *testing.T) {
	client := NewAPIClient("http://localhost:8080")

	tests := []struct {
		name       string
		path       string
		expectedCode int
	}{
		{"Non-existent endpoint", "/api/v1/nonexistent", http.StatusNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, err := client.Get(test.path)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != test.expectedCode {
				t.Errorf("Expected %d, got %d", test.expectedCode, resp.StatusCode)
			}
		})
	}
}

func TestDataConsistency(t *testing.T) {
	client := NewAPIClient("http://localhost:8080")

	// Get metrics multiple times
	var results []map[string]interface{}
	for i := 0; i < 3; i++ {
		resp, err := client.Get("/api/v1/metrics")
		if err != nil {
			t.Fatalf("Request %d failed: %v", i, err)
		}

		var data map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&data)
		resp.Body.Close()

		results = append(results, data)
	}

	// Verify metrics exist in all results
	for i, result := range results {
		if result == nil {
			t.Errorf("Result %d is nil", i)
		}
	}
}
