-- Widget data persistence table
CREATE TABLE IF NOT EXISTS widget_data (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    widget_id TEXT NOT NULL,
    widget_type TEXT NOT NULL,
    data_key TEXT NOT NULL,
    data_value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(widget_id, data_key)
);

CREATE INDEX IF NOT EXISTS idx_widget_data_lookup ON widget_data(widget_id, data_key);
CREATE INDEX IF NOT EXISTS idx_widget_data_updated ON widget_data(updated_at);

-- Historical data tracking table
CREATE TABLE IF NOT EXISTS widget_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    widget_id TEXT NOT NULL,
    widget_type TEXT NOT NULL,
    metric_name TEXT NOT NULL,
    metric_value REAL NOT NULL,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_history_widget_time ON widget_history(widget_id, recorded_at);
CREATE INDEX IF NOT EXISTS idx_history_metric ON widget_history(metric_name, recorded_at);

-- Activity/audit logs table
CREATE TABLE IF NOT EXISTS activity_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_type TEXT NOT NULL,
    widget_id TEXT,
    user_id TEXT,
    details TEXT,
    ip_address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_activity_time ON activity_log(created_at);
CREATE INDEX IF NOT EXISTS idx_activity_event ON activity_log(event_type);

-- Custom alerts table
CREATE TABLE IF NOT EXISTS alerts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    widget_id TEXT NOT NULL,
    condition_type TEXT NOT NULL,
    condition_value TEXT NOT NULL,
    notification_type TEXT NOT NULL,
    notification_target TEXT NOT NULL,
    enabled BOOLEAN DEFAULT 1,
    last_triggered TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_alerts_widget ON alerts(widget_id);
CREATE INDEX IF NOT EXISTS idx_alerts_enabled ON alerts(enabled);

-- Dashboard profiles/presets table
CREATE TABLE IF NOT EXISTS dashboard_profiles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    config_json TEXT NOT NULL,
    is_active BOOLEAN DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_profiles_active ON dashboard_profiles(is_active);

-- API rate limiting table
CREATE TABLE IF NOT EXISTS rate_limits (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    identifier TEXT NOT NULL,
    endpoint TEXT NOT NULL,
    request_count INTEGER DEFAULT 0,
    window_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(identifier, endpoint, window_start)
);

CREATE INDEX IF NOT EXISTS idx_rate_limits_cleanup ON rate_limits(window_start);
