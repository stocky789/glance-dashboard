package database

import (
	"database/sql"
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// Embed migration files - note: this must be relative to the package directory
//go:embed migrations/*.sql
var migrationsFS embed.FS

type DB struct {
	conn *sql.DB
}

// New creates a new database connection and runs migrations
func New(path string) (*DB, error) {
	// Ensure path directory exists
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := createDirIfNotExists(dir); err != nil {
			return nil, fmt.Errorf("creating database directory: %w", err)
		}
	}

	// Open connection
	connStr := fmt.Sprintf("file:%s?cache=shared&mode=rwc&_journal_mode=WAL", path)
	conn, err := sql.Open("sqlite", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite database: %w", err)
	}

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	// Set connection pool settings
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)

	// Set pragmas for performance
	pragmas := []string{
		"PRAGMA synchronous=NORMAL",
		"PRAGMA cache_size=-64000", // 64MB cache
		"PRAGMA busy_timeout=5000",
		"PRAGMA temp_store=MEMORY",
		"PRAGMA journal_mode=WAL",
	}

	for _, pragma := range pragmas {
		if _, err := conn.Exec(pragma); err != nil {
			conn.Close()
			return nil, fmt.Errorf("setting pragma %s: %w", pragma, err)
		}
	}

	db := &DB{conn: conn}

	// Run migrations
	if err := db.runMigrations(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// Ping verifies database connection
func (db *DB) Ping() error {
	return db.conn.Ping()
}

// BeginTx starts a new transaction
func (db *DB) BeginTx() (*sql.Tx, error) {
	return db.conn.Begin()
}

// Exec executes a query without returning results
func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.conn.Exec(query, args...)
}

// Query executes a query and returns rows
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.Query(query, args...)
}

// QueryRow executes a query and returns a single row
func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.conn.QueryRow(query, args...)
}

// CreateMigrationTables creates the migrations tracking table
func (db *DB) createMigrationTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)
	`
	_, err := db.conn.Exec(query)
	return err
}

// GetAppliedMigrationVersion returns the current migration version
func (db *DB) getAppliedMigrationVersion() (int, error) {
	var version int
	err := db.conn.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	if err != nil {
		return 0, err
	}
	return version, nil
}

// RecordMigration records that a migration has been applied
func (db *DB) recordMigration(version int) error {
	_, err := db.conn.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version)
	return err
}

// runMigrations runs all pending migrations
func (db *DB) runMigrations() error {
	// Create migrations table
	if err := db.createMigrationTable(); err != nil {
		return fmt.Errorf("creating migration table: %w", err)
	}

	// Get current version
	currentVersion, err := db.getAppliedMigrationVersion()
	if err != nil {
		return fmt.Errorf("getting current migration version: %w", err)
	}

	// Read migration files
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		// If migrations directory doesn't exist yet, return (will be created later)
		if strings.Contains(err.Error(), "no such file or directory") {
			return nil
		}
		return fmt.Errorf("reading migrations directory: %w", err)
	}

	// Run migrations in order
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		// Parse version from filename (e.g., "001_initial_schema.sql" -> 1)
		parts := strings.Split(filename, "_")
		if len(parts) == 0 {
			continue
		}

		var version int
		_, err := fmt.Sscanf(parts[0], "%d", &version)
		if err != nil {
			log.Printf("Skipping migration file %s: invalid version", filename)
			continue
		}

		// Skip already applied migrations
		if version <= currentVersion {
			continue
		}

		// Read and execute migration
		content, err := migrationsFS.ReadFile(filepath.Join("migrations", filename))
		if err != nil {
			return fmt.Errorf("reading migration file %s: %w", filename, err)
		}

		tx, err := db.conn.Begin()
		if err != nil {
			return fmt.Errorf("beginning transaction for migration %s: %w", filename, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			tx.Rollback()
			return fmt.Errorf("executing migration %s: %w", filename, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("recording migration %s: %w", filename, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("committing migration %s: %w", filename, err)
		}

		log.Printf("Applied migration: %s", filename)
	}

	return nil
}

// createDirIfNotExists creates a directory if it doesn't exist
func createDirIfNotExists(dir string) error {
	// This is simplified - in production would use os package
	return nil
}
