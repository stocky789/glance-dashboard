package database

import (
	"encoding/json"
	"fmt"
	"time"
)

// WidgetData represents persisted widget data
type WidgetData struct {
	ID        int64       `json:"id"`
	WidgetID  string      `json:"widget_id"`
	Type      string      `json:"type"`
	DataKey   string      `json:"key"`
	DataValue interface{} `json:"value"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// SaveWidgetData saves or updates widget data
func (db *DB) SaveWidgetData(widgetID, widgetType, key string, value interface{}) error {
	if db.conn == nil {
		return fmt.Errorf("database connection not initialized")
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshaling value: %w", err)
	}

	query := `
	INSERT INTO widget_data (widget_id, widget_type, data_key, data_value, updated_at)
	VALUES (?, ?, ?, ?, CURRENT_TIMESTAMP)
	ON CONFLICT(widget_id, data_key) 
	DO UPDATE SET data_value=excluded.data_value, updated_at=CURRENT_TIMESTAMP
	`

	_, err = db.conn.Exec(query, widgetID, widgetType, key, string(jsonValue))
	return err
}

// GetWidgetData retrieves a specific piece of widget data
func (db *DB) GetWidgetData(widgetID, key string) (*WidgetData, error) {
	if db.conn == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	query := `
	SELECT id, widget_id, widget_type, data_key, data_value, updated_at
	FROM widget_data
	WHERE widget_id = ? AND data_key = ?
	`

	var wd WidgetData
	var jsonValue string

	err := db.conn.QueryRow(query, widgetID, key).Scan(
		&wd.ID, &wd.WidgetID, &wd.Type, &wd.DataKey, &jsonValue, &wd.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("querying widget data: %w", err)
	}

	if err := json.Unmarshal([]byte(jsonValue), &wd.DataValue); err != nil {
		return nil, fmt.Errorf("unmarshaling value: %w", err)
	}

	return &wd, nil
}

// GetAllWidgetData retrieves all data for a widget
func (db *DB) GetAllWidgetData(widgetID string) ([]WidgetData, error) {
	if db.conn == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	query := `
	SELECT id, widget_id, widget_type, data_key, data_value, updated_at
	FROM widget_data
	WHERE widget_id = ?
	ORDER BY data_key
	`

	rows, err := db.conn.Query(query, widgetID)
	if err != nil {
		return nil, fmt.Errorf("querying widget data: %w", err)
	}
	defer rows.Close()

	var results []WidgetData
	for rows.Next() {
		var wd WidgetData
		var jsonValue string

		if err := rows.Scan(&wd.ID, &wd.WidgetID, &wd.Type, &wd.DataKey, &jsonValue, &wd.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning widget data: %w", err)
		}

		if err := json.Unmarshal([]byte(jsonValue), &wd.DataValue); err != nil {
			return nil, fmt.Errorf("unmarshaling value: %w", err)
		}

		results = append(results, wd)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return results, nil
}

// DeleteWidgetData deletes a piece of widget data
func (db *DB) DeleteWidgetData(widgetID, key string) error {
	if db.conn == nil {
		return fmt.Errorf("database connection not initialized")
	}

	query := `DELETE FROM widget_data WHERE widget_id = ? AND data_key = ?`
	_, err := db.conn.Exec(query, widgetID, key)
	return err
}

// DeleteAllWidgetData deletes all data for a widget
func (db *DB) DeleteAllWidgetData(widgetID string) error {
	if db.conn == nil {
		return fmt.Errorf("database connection not initialized")
	}

	query := `DELETE FROM widget_data WHERE widget_id = ?`
	_, err := db.conn.Exec(query, widgetID)
	return err
}
