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
}

func (w *metricsWidget) renderMetrics(m metricsData) {
	tmpl := template.Must(template.New("metrics").Parse(`
		<div class="list list-gap-12">
			<div class="list-horizontal list-gap-10">
				<span class="size-h4 color-primary">{{.SystemMetrics.Memory}}</span>
				<span>MB Memory</span>
			</div>
			<div class="list-horizontal list-gap-10">
				<span class="size-h4 color-primary">{{.SystemMetrics.Goroutines}}</span>
				<span>Goroutines</span>
			</div>
			<div class="list-horizontal list-gap-10">
				<span class="size-h4 color-primary">{{.SystemMetrics.Uptime}}</span>
				<span>Uptime</span>
			</div>
			<div class="list-horizontal list-gap-10">
				<span class="size-h4 color-primary">{{.APIMetrics.TotalRequests}}</span>
				<span>API Requests</span>
			</div>
			<div class="list-horizontal list-gap-10">
				<span class="size-h4 color-primary">{{.APIMetrics.AverageLatency | printf "%.1f"}}</span>
				<span>ms Latency</span>
			</div>
			<div class="list-horizontal list-gap-10">
				<span class="size-h4 color-primary">{{.APIMetrics.RateLimited}}</span>
				<span>Rate Limited</span>
			</div>
		</div>
	`))
	w.templateBuffer.Reset()
	tmpl.Execute(&w.templateBuffer, m)
	w.ContentAvailable = true
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
