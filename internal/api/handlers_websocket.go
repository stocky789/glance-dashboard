package api

import (
	"log"
	"net/http"

	wsinternal "github.com/glanceapp/glance/internal/websocket"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking based on configuration
		return true
	},
}

// handleWebSocket handles WebSocket upgrade requests
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	if s.wsHub == nil {
		http.Error(w, "WebSocket service not available", http.StatusServiceUnavailable)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := wsinternal.NewClient(s.wsHub, conn)
	s.wsHub.RegisterClient(client)

	go client.WritePump()
	go client.ReadPump()
}

// handleMetrics returns system and widget performance metrics
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if s.metricsCollector == nil {
		http.Error(w, "Metrics service not available", http.StatusServiceUnavailable)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := map[string]interface{}{
		"system":  s.metricsCollector.GetSystemMetrics(),
		"widgets": s.metricsCollector.GetAllMetrics(),
	}

	w.Header().Set("Content-Type", "application/json")
	encodeJSON(w, metrics)
}

// handleMetricsWidget returns metrics for a specific widget
func (s *Server) handleMetricsWidget(w http.ResponseWriter, r *http.Request, widgetID string) {
	if s.metricsCollector == nil {
		http.Error(w, "Metrics service not available", http.StatusServiceUnavailable)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := s.metricsCollector.GetMetrics(widgetID)
	if metrics == nil {
		http.Error(w, "Widget not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encodeJSON(w, metrics)
}
