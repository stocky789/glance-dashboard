package database

import (
	"time"
)

type HistoricalDataPoint struct {
	ID         int64
	WidgetID   string
	WidgetType string
	MetricName string
	Value      float64
	RecordedAt time.Time
}

func (db *DB) RecordMetric(widgetID, widgetType, metricName string, value float64) error {
	query := `
		INSERT INTO widget_history (widget_id, widget_type, metric_name, metric_value, recorded_at)
		VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
	`
	_, err := db.conn.Exec(query, widgetID, widgetType, metricName, value)
	return err
}

func (db *DB) GetHistory(widgetID, metricName string, since time.Time) ([]HistoricalDataPoint, error) {
	query := `
		SELECT id, widget_id, widget_type, metric_name, metric_value, recorded_at
		FROM widget_history
		WHERE widget_id = ? AND metric_name = ? AND recorded_at >= ?
		ORDER BY recorded_at ASC
	`
	rows, err := db.conn.Query(query, widgetID, metricName, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []HistoricalDataPoint
	for rows.Next() {
		var point HistoricalDataPoint
		err := rows.Scan(&point.ID, &point.WidgetID, &point.WidgetType, &point.MetricName, &point.Value, &point.RecordedAt)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}
	return points, rows.Err()
}

func (db *DB) CleanupOldHistory(retentionDays int) error {
	query := `
		DELETE FROM widget_history
		WHERE recorded_at < datetime('now', ? || ' days')
	`
	_, err := db.conn.Exec(query, -retentionDays)
	return err
}
