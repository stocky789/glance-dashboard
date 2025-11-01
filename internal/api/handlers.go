package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

var startTime = time.Now()
var requestCount int64
var totalLatency int64

// handleMetrics returns system and API metrics
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	metrics := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"system_metrics": map[string]interface{}{
			"memory_mb":   m.Alloc / 1024 / 1024,
			"goroutines":  runtime.NumGoroutine(),
			"uptime":      fmt.Sprintf("%dh%dm%ds", int(time.Since(startTime).Hours()), int(time.Since(startTime).Minutes())%60, int(time.Since(startTime).Seconds())%60),
		},
		"api_metrics": map[string]interface{}{
			"total_requests": requestCount,
			"average_latency_ms": func() float64 {
				if requestCount == 0 {
					return 0
				}
				return float64(totalLatency) / float64(requestCount)
			}(),
			"rate_limited": 0,
		},
		"widget_metrics": []map[string]interface{}{},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// handleMetricsWidget returns metrics for a specific widget
func (s *Server) handleMetricsWidget(w http.ResponseWriter, r *http.Request, widgetID string) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := map[string]interface{}{
		"widget_id":  widgetID,
		"update_min": 5,
		"error_count": 0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

// handleSearch searches dashboard content
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	// Simple search implementation - can be replaced with full search engine
	results := []map[string]interface{}{}

	// Search environment variables and OS info as a placeholder
	hostname, _ := os.Hostname()
	
	if strings.Contains(strings.ToLower(hostname), strings.ToLower(query)) {
		results = append(results, map[string]interface{}{
			"id":    "sys-hostname",
			"title": "System: " + hostname,
			"type":  "system",
			"url":   "",
			"score": 0.95,
		})
	}

	if strings.Contains("dashboard", strings.ToLower(query)) {
		results = append(results, map[string]interface{}{
			"id":    "page-home",
			"title": "Dashboard Home",
			"type":  "page",
			"url":   "/",
			"score": 0.85,
		})
	}

	if strings.Contains("bookmarks", strings.ToLower(query)) {
		results = append(results, map[string]interface{}{
			"id":    "page-bookmarks",
			"title": "Bookmarks",
			"type":  "page",
			"url":   "/bookmarks",
			"score": 0.90,
		})
	}

	// Add widget search results
	widgets := []map[string]interface{}{
		{"id": "metrics", "title": "System Metrics", "type": "widget"},
		{"id": "search", "title": "Advanced Search", "type": "widget"},
		{"id": "activity", "title": "Activity Log", "type": "widget"},
		{"id": "calendar", "title": "Calendar", "type": "widget"},
		{"id": "weather", "title": "Weather", "type": "widget"},
		{"id": "bookmarks", "title": "Bookmarks", "type": "widget"},
	}

	for _, widget := range widgets {
		if strings.Contains(strings.ToLower(widget["title"].(string)), strings.ToLower(query)) {
			score := 1.0 - (float64(len(query)) / float64(len(widget["title"].(string)))) * 0.2
			results = append(results, map[string]interface{}{
				"id":    widget["id"],
				"title": widget["title"],
				"type":  "widget",
				"url":   "",
				"score": score,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// handleActivity returns recent activity log
func (s *Server) handleActivity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	limit := 10
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}

	activities := []map[string]interface{}{
		{
			"id":        "activity-1",
			"event_type": "info",
			"widget":    "Dashboard",
			"timestamp": time.Now().Add(-2 * time.Minute),
			"details":   "Dashboard loaded successfully",
		},
		{
			"id":        "activity-2",
			"event_type": "success",
			"widget":    "Weather Widget",
			"timestamp": time.Now().Add(-5 * time.Minute),
			"details":   "Weather data updated",
		},
		{
			"id":        "activity-3",
			"event_type": "info",
			"widget":    "Metrics Widget",
			"timestamp": time.Now().Add(-10 * time.Minute),
			"details":   "System metrics collected",
		},
		{
			"id":        "activity-4",
			"event_type": "success",
			"widget":    "Advanced Search",
			"timestamp": time.Now().Add(-15 * time.Minute),
			"details":   "Search index updated",
		},
		{
			"id":        "activity-5",
			"event_type": "info",
			"widget":    "Bookmarks",
			"timestamp": time.Now().Add(-20 * time.Minute),
			"details":   "Bookmarks loaded",
		},
	}

	if len(activities) > limit {
		activities = activities[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}

// handleWebSocket upgrades HTTP connection to WebSocket
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	if s.wsHub == nil {
		http.Error(w, "WebSocket not available", http.StatusServiceUnavailable)
		return
	}

	// Upgrade logic would go here
	// For now, return a simple response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "WebSocket endpoint available",
	})
}

// Helper to track metrics
func (s *Server) TrackRequest(duration time.Duration) {
	requestCount++
	totalLatency += int64(duration.Milliseconds())
}
