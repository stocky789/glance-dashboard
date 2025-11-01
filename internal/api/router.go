package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/glanceapp/glance/internal/database"
	"github.com/glanceapp/glance/internal/websocket"
)

// MetricsCollector interface for metrics collection
type MetricsCollector interface {
	GetSystemMetrics() map[string]interface{}
	GetAllMetrics() map[string]interface{}
	GetMetrics(widgetID string) interface{}
}

// Server represents the API server with all handlers and middleware
type Server struct {
	db               *database.DB
	mux              *http.ServeMux
	config           *Config
	wsHub            *websocket.Hub
	metricsCollector MetricsCollector
}

// Config holds API server configuration
type Config struct {
	RateLimitEnabled bool
	RateLimitRPM     int
	CORSEnabled      bool
	CORSOrigins      []string
}

// NewServer creates a new API server instance
func NewServer(db *database.DB, config *Config) *Server {
	s := &Server{
		db:     db,
		mux:    http.NewServeMux(),
		config: config,
	}

	s.setupRoutes()
	return s
}

// SetWebSocketHub sets the WebSocket hub for real-time communication
func (s *Server) SetWebSocketHub(hub *websocket.Hub) {
	s.wsHub = hub
}

// SetMetricsCollector sets the metrics collector for performance monitoring
func (s *Server) SetMetricsCollector(collector MetricsCollector) {
	s.metricsCollector = collector
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() {
	// Health check
	s.mux.HandleFunc("/api/health", s.handleHealth)

	// WebSocket endpoint
	s.mux.HandleFunc("/api/ws", s.handleWebSocket)

	// Metrics endpoints
	s.mux.HandleFunc("/api/v1/metrics", s.handleMetrics)
	s.mux.HandleFunc("/api/v1/metrics/widgets/", s.routeMetricsEndpoints)

	// Widget data endpoints
	s.mux.HandleFunc("/api/v1/widgets/", s.routeWidgetEndpoints)
}

// routeWidgetEndpoints routes widget-related endpoints
func (s *Server) routeWidgetEndpoints(w http.ResponseWriter, r *http.Request) {
	// Parse path: /api/v1/widgets/{id} or /api/v1/widgets/{id}/data or /api/v1/widgets/{id}/data/{key}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/widgets/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 {
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}

	widgetID := parts[0]

	if len(parts) == 1 {
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}

	if parts[1] != "data" {
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}

	// GET /api/v1/widgets/{id}/data or POST /api/v1/widgets/{id}/data
	if len(parts) == 2 {
		s.handleWidgetData(w, r, widgetID)
		return
	}

	// GET /api/v1/widgets/{id}/data/{key} or DELETE /api/v1/widgets/{id}/data/{key}
	if len(parts) == 3 {
		key := parts[2]
		s.handleWidgetDataKey(w, r, widgetID, key)
		return
	}

	http.Error(w, "Invalid endpoint", http.StatusBadRequest)
}

// routeMetricsEndpoints routes metrics-related endpoints
func (s *Server) routeMetricsEndpoints(w http.ResponseWriter, r *http.Request) {
	// Parse path: /api/v1/metrics/widgets/{id}
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/metrics/widgets/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}

	widgetID := parts[0]
	s.handleMetricsWidget(w, r, widgetID)
}

// handleWidgetData handles GET and POST for all widget data
func (s *Server) handleWidgetData(w http.ResponseWriter, r *http.Request, widgetID string) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetAllWidgetData(w, r, widgetID)
	case http.MethodPost:
		s.handleSaveWidgetData(w, r, widgetID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleWidgetDataKey handles GET and DELETE for specific widget data
func (s *Server) handleWidgetDataKey(w http.ResponseWriter, r *http.Request, widgetID, key string) {
	switch r.Method {
	case http.MethodGet:
		s.handleGetWidgetData(w, r, widgetID, key)
	case http.MethodDelete:
		s.handleDeleteWidgetData(w, r, widgetID, key)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handleGetAllWidgetData retrieves all data for a widget
func (s *Server) handleGetAllWidgetData(w http.ResponseWriter, r *http.Request, widgetID string) {
	data, err := s.db.GetAllWidgetData(widgetID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve widget data: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// handleSaveWidgetData saves a piece of widget data
func (s *Server) handleSaveWidgetData(w http.ResponseWriter, r *http.Request, widgetID string) {
	var payload struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
		Type  string      `json:"type,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if payload.Key == "" {
		http.Error(w, "Missing 'key' field", http.StatusBadRequest)
		return
	}

	// Default type if not provided
	widgetType := payload.Type
	if widgetType == "" {
		widgetType = "unknown"
	}

	if err := s.db.SaveWidgetData(widgetID, widgetType, payload.Key, payload.Value); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save widget data: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "created",
		"widget_id": widgetID,
		"key": payload.Key,
	})
}

// handleGetWidgetData retrieves a specific piece of widget data
func (s *Server) handleGetWidgetData(w http.ResponseWriter, r *http.Request, widgetID, key string) {
	data, err := s.db.GetWidgetData(widgetID, key)
	if err != nil {
		http.Error(w, fmt.Sprintf("Widget data not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// handleDeleteWidgetData deletes a specific piece of widget data
func (s *Server) handleDeleteWidgetData(w http.ResponseWriter, r *http.Request, widgetID, key string) {
	if err := s.db.DeleteWidgetData(widgetID, key); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete widget data: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleHealth provides health check information
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := struct {
		Status   string `json:"status"`
		Database string `json:"database"`
		Version  string `json:"version"`
	}{
		Status:   "ok",
		Database: "connected",
		Version:  "0.9.0",
	}

	// Check database connection
	if s.db != nil {
		if err := s.db.Ping(); err != nil {
			health.Status = "degraded"
			health.Database = "disconnected"
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	encodeJSON(w, health)
}

// encodeJSON is a helper function to encode JSON responses
func encodeJSON(w http.ResponseWriter, data interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// ServeHTTP implements http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Apply middleware chain
	handler := http.Handler(s.mux)

	if s.config.CORSEnabled {
		handler = s.corsMiddleware(handler)
	}

	if s.config.RateLimitEnabled {
		handler = s.rateLimitMiddleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// corsMiddleware adds CORS headers to responses
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range s.config.CORSOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "3600")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// rateLimitMiddleware applies rate limiting to requests
func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For now, implement a simple pass-through
		// Rate limiting will be enhanced in Phase 2
		next.ServeHTTP(w, r)
	})
}
