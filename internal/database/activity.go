package database

import (
	"encoding/json"
	"strings"
	"time"
)

type ActivityLog struct {
	ID        int64
	EventType string
	WidgetID  string
	UserID    string
	Details   map[string]interface{}
	IPAddress string
	CreatedAt time.Time
}

func (db *DB) LogActivity(eventType, widgetID, userID, ipAddress string, details map[string]interface{}) error {
	detailsJSON, _ := json.Marshal(details)
	query := `
		INSERT INTO activity_log (event_type, widget_id, user_id, details, ip_address, created_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`
	_, err := db.conn.Exec(query, eventType, widgetID, userID, string(detailsJSON), ipAddress)
	return err
}

func (db *DB) GetActivityLog(since time.Time, eventTypes []string, limit int) ([]ActivityLog, error) {
	query := `
		SELECT id, event_type, widget_id, user_id, details, ip_address, created_at
		FROM activity_log
		WHERE created_at >= ?
	`
	args := []interface{}{since}

	if len(eventTypes) > 0 {
		placeholders := make([]string, len(eventTypes))
		for i, et := range eventTypes {
			placeholders[i] = "?"
			args = append(args, et)
		}
		query += " AND event_type IN (" + strings.Join(placeholders, ",") + ")"
	}

	query += " ORDER BY created_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []ActivityLog
	for rows.Next() {
		var log ActivityLog
		var detailsJSON string
		err := rows.Scan(&log.ID, &log.EventType, &log.WidgetID, &log.UserID, &detailsJSON, &log.IPAddress, &log.CreatedAt)
		if err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(detailsJSON), &log.Details)
		logs = append(logs, log)
	}

	return logs, rows.Err()
}
