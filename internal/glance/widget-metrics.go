package glance

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"time"
)

type metricsWidget struct {
	widgetBase
}

type metricsData struct {
	Timestamp     int64            `json:"timestamp"`
	SystemMetrics systemMetrics    `json:"system_metrics"`
	WidgetMetrics []widgetMetric   `json:"widget_metrics"`
	APIMetrics    apiMetrics       `json:"api_metrics"`
}

type systemMetrics struct {
	Memory     uint64 `json:"memory_mb"`
	Goroutines int    `json:"goroutines"`
	Uptime     string `json:"uptime"`
}

type widgetMetric struct {
	WidgetID  string `json:"widget_id"`
	UpdateMin int    `json:"update_min"`
	ErrorCount int   `json:"error_count"`
}

type apiMetrics struct {
	TotalRequests   int64   `json:"total_requests"`
	AverageLatency  float64 `json:"average_latency_ms"`
	RateLimited     int64   `json:"rate_limited"`
}

func (w *metricsWidget) initialize() error {
	w.cacheDuration = 5 * time.Second
	w.cacheType = cacheTypeDuration
	w.ContentAvailable = true
	return nil
}

func (w *metricsWidget) update(ctx context.Context) {
	resp, err := http.Get("http://localhost:8080/api/v1/metrics")
	if err != nil {
		w.Error = err
		return
	}
	defer resp.Body.Close()

	var metrics metricsData
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		w.Error = err
		return
	}

	w.renderMetrics(metrics)
	w.ContentAvailable = true
}

func (w *metricsWidget) renderMetrics(m metricsData) {
	tmpl := template.Must(template.New("metrics").Parse(`
		<div class="metrics-dashboard">
			<div class="metrics-section">
				<h3>System</h3>
				<div class="metric-item">
					<span>Memory:</span>
					<strong>{{.SystemMetrics.Memory}} MB</strong>
				</div>
				<div class="metric-item">
					<span>Goroutines:</span>
					<strong>{{.SystemMetrics.Goroutines}}</strong>
				</div>
				<div class="metric-item">
					<span>Uptime:</span>
					<strong>{{.SystemMetrics.Uptime}}</strong>
				</div>
			</div>
			<div class="metrics-section">
				<h3>API</h3>
				<div class="metric-item">
					<span>Requests:</span>
					<strong>{{.APIMetrics.TotalRequests}}</strong>
				</div>
				<div class="metric-item">
					<span>Avg Latency:</span>
					<strong>{{.APIMetrics.AverageLatency | printf "%.1f"}} ms</strong>
				</div>
				<div class="metric-item">
					<span>Rate Limited:</span>
					<strong>{{.APIMetrics.RateLimited}}</strong>
				</div>
			</div>
		</div>
		<style>
			.metrics-dashboard {
				display: grid;
				grid-template-columns: 1fr 1fr;
				gap: 1rem;
				padding: 1rem;
			}
			.metrics-section {
				background: rgba(0,0,0,0.1);
				border-radius: 8px;
				padding: 0.75rem;
			}
			.metrics-section h3 {
				margin: 0 0 0.5rem 0;
				font-size: 0.9rem;
				opacity: 0.7;
			}
			.metric-item {
				display: flex;
				justify-content: space-between;
				padding: 0.25rem 0;
				font-size: 0.85rem;
			}
			.metric-item strong {
				color: var(--primary-color, #00ff00);
				font-weight: bold;
			}
		</style>
	`))

	w.templateBuffer.Reset()
	tmpl.Execute(&w.templateBuffer, m)
}

func (w *metricsWidget) Render() template.HTML {
	return template.HTML(w.templateBuffer.String())
}

func (w *metricsWidget) GetType() string {
	return "metrics"
}

func (w *metricsWidget) handleRequest(wr http.ResponseWriter, r *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(map[string]string{"status": "ok"})
}

func (w *metricsWidget) setHideHeader(b bool) {
	w.HideHeader = b
}
