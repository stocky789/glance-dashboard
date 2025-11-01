package database

import "database/sql"

const schema = `
CREATE TABLE IF NOT EXISTS widgets (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	type TEXT NOT NULL,
	config TEXT,
	enabled BOOLEAN DEFAULT 1,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS metrics (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	widget_id TEXT NOT NULL,
	update_count INTEGER DEFAULT 0,
	error_count INTEGER DEFAULT 0,
	avg_update_duration REAL DEFAULT 0,
	last_update TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (widget_id) REFERENCES widgets(id)
);

CREATE TABLE IF NOT EXISTS history (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	widget_id TEXT NOT NULL,
	data TEXT NOT NULL,
	timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (widget_id) REFERENCES widgets(id)
);

CREATE TABLE IF NOT EXISTS activity_log (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	event_type TEXT NOT NULL,
	widget_id TEXT,
	user_id TEXT,
	details TEXT,
	ip_address TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (widget_id) REFERENCES widgets(id)
);

CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
	user_id TEXT,
	data TEXT,
	expires_at TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_metrics_widget ON metrics(widget_id);
CREATE INDEX IF NOT EXISTS idx_history_widget ON history(widget_id);
CREATE INDEX IF NOT EXISTS idx_history_timestamp ON history(timestamp);
CREATE INDEX IF NOT EXISTS idx_activity_event_type ON activity_log(event_type);
CREATE INDEX IF NOT EXISTS idx_activity_created_at ON activity_log(created_at);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);
`

func (db *DB) InitializeSchema() error {
	_, err := db.conn.Exec(schema)
	return err
}

func (db *DB) CreateWidget(id, name, widgetType, config string) error {
	query := `
		INSERT INTO widgets (id, name, type, config)
		VALUES (?, ?, ?, ?)
	`
	_, err := db.conn.Exec(query, id, name, widgetType, config)
	return err
}

func (db *DB) GetWidget(id string) (map[string]interface{}, error) {
	query := `
		SELECT id, name, type, config, enabled, created_at, updated_at
		FROM widgets
		WHERE id = ?
	`
	row := db.conn.QueryRow(query, id)

	var widgetID, name, widgetType string
	var widgetConfig sql.NullString
	var enabled bool
	var createdAt, updatedAt sql.NullTime

	err := row.Scan(&widgetID, &name, &widgetType, &widgetConfig, &enabled, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	widget := map[string]interface{}{
		"id":         widgetID,
		"name":       name,
		"type":       widgetType,
		"config":     widgetConfig.String,
		"enabled":    enabled,
		"created_at": createdAt.Time,
		"updated_at": updatedAt.Time,
	}

	return widget, nil
}

func (db *DB) ListWidgets() ([]map[string]interface{}, error) {
	query := `
		SELECT id, name, type, config, enabled, created_at, updated_at
		FROM widgets
		ORDER BY created_at DESC
	`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var widgets []map[string]interface{}
	return widgets, nil
}

func (db *DB) UpdateWidget(id, name, widgetType, config string) error {
	query := `
		UPDATE widgets
		SET name = ?, type = ?, config = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`
	_, err := db.conn.Exec(query, name, widgetType, config, id)
	return err
}

func (db *DB) DeleteWidget(id string) error {
	query := `DELETE FROM widgets WHERE id = ?`
	_, err := db.conn.Exec(query, id)
	return err
}

func (db *DB) DisableWidget(id string) error {
	query := `UPDATE widgets SET enabled = 0 WHERE id = ?`
	_, err := db.conn.Exec(query, id)
	return err
}

func (db *DB) EnableWidget(id string) error {
	query := `UPDATE widgets SET enabled = 1 WHERE id = ?`
	_, err := db.conn.Exec(query, id)
	return err
}

func (db *DB) GetEnabledWidgets() ([]string, error) {
	query := `SELECT id FROM widgets WHERE enabled = 1`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}

func (db *DB) CleanupOldData(daysToKeep int) error {
	query := `
		DELETE FROM history
		WHERE timestamp < datetime('now', '-' || ? || ' days')
	`
	_, err := db.conn.Exec(query, daysToKeep)
	return err
}
