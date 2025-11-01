package glance

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type activityLogWidget struct {
	widgetBase
	Limit int `yaml:"limit"`
}

type activity struct {
	ID        string    `json:"id"`
	EventType string    `json:"event_type"`
	Widget    string    `json:"widget"`
	Timestamp time.Time `json:"timestamp"`
	Details   string    `json:"details"`
}

func (w *activityLogWidget) initialize() error {
	if w.Limit == 0 {
		w.Limit = 10
	}
	w.Title = "Activity Log"
	w.cacheDuration = 10 * time.Second
	w.cacheType = cacheTypeDuration
	return nil
}

func (w *activityLogWidget) update(ctx context.Context) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/api/v1/activity?limit=%d", w.Limit))
	if err != nil {
		w.Error = err
		return
	}
	defer resp.Body.Close()

	var activities []activity
	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		w.Error = err
		return
	}

	w.renderActivity(activities)
	w.ContentAvailable = true
}

func (w *activityLogWidget) renderActivity(activities []activity) {
	tmpl := template.Must(template.New("activity").Parse(`
		<div class="activity-log">
			{{range .Activities}}
			<div class="activity-item {{.EventType}}">
				<div class="activity-header">
					<span class="activity-type">{{.EventType}}</span>
					<span class="activity-time" title="{{.Timestamp.Format "Mon Jan 02 15:04:05"}}">
						{{.Timestamp.Format "15:04"}}
					</span>
				</div>
				<div class="activity-widget">{{.Widget}}</div>
				<div class="activity-details">{{.Details}}</div>
			</div>
			{{end}}
		</div>
		<style>
			.activity-log {
				padding: 0.5rem;
				max-height: 400px;
				overflow-y: auto;
			}
			.activity-item {
				padding: 0.75rem;
				margin: 0.5rem 0;
				border-left: 3px solid rgba(255,255,255,0.2);
				border-radius: 4px;
				background: rgba(0,0,0,0.2);
				font-size: 0.85rem;
			}
			.activity-item.error {
				border-left-color: #ff6b6b;
				background: rgba(255,107,107,0.1);
			}
			.activity-item.success {
				border-left-color: #51cf66;
				background: rgba(81,207,102,0.1);
			}
			.activity-item.warning {
				border-left-color: #ffd43b;
				background: rgba(255,212,59,0.1);
			}
			.activity-item.info {
				border-left-color: #4dabf7;
				background: rgba(77,171,247,0.1);
			}
			.activity-header {
				display: flex;
				justify-content: space-between;
				align-items: center;
				margin-bottom: 0.25rem;
			}
			.activity-type {
				font-weight: bold;
				text-transform: uppercase;
				font-size: 0.75rem;
			}
			.activity-item.error .activity-type {
				color: #ff6b6b;
			}
			.activity-item.success .activity-type {
				color: #51cf66;
			}
			.activity-item.warning .activity-type {
				color: #ffd43b;
			}
			.activity-item.info .activity-type {
				color: #4dabf7;
			}
			.activity-time {
				font-size: 0.75rem;
				opacity: 0.7;
			}
			.activity-widget {
				font-weight: 600;
				margin: 0.25rem 0;
				color: var(--primary-color, #00ff00);
			}
			.activity-details {
				opacity: 0.8;
				font-size: 0.8rem;
				margin-top: 0.25rem;
			}
		</style>
	`))

	w.templateBuffer.Reset()
	tmpl.Execute(&w.templateBuffer, map[string]interface{}{
		"Activities": activities,
	})
}

func (w *activityLogWidget) Render() template.HTML {
	return template.HTML(w.templateBuffer.String())
}

func (w *activityLogWidget) GetType() string {
	return "activity-log"
}

func (w *activityLogWidget) handleRequest(wr http.ResponseWriter, r *http.Request) {
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(map[string]string{"status": "ok"})
}

func (w *activityLogWidget) setHideHeader(b bool) {
	w.HideHeader = b
}
