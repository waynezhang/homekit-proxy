package actionlog

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ActionLog struct {
	db *sql.DB
}

type LogEntry struct {
	ID                 int64     `json:"id"`
	EntityID           string    `json:"entity_id"`
	CharacteristicType string    `json:"characteristic_type"`
	NewValue           string    `json:"new_value"`
	Timestamp          time.Time `json:"timestamp"`
}

func New(dbPath string) (*ActionLog, error) {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	al := &ActionLog{db: db}
	if err := al.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return al, nil
}

func (al *ActionLog) initSchema() error {
	query := `
	CREATE TABLE IF NOT EXISTS action_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		entity_id TEXT NOT NULL,
		characteristic_type TEXT NOT NULL,
		new_value TEXT,
		timestamp DATETIME NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_entity_id ON action_logs(entity_id);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON action_logs(timestamp);
	CREATE INDEX IF NOT EXISTS idx_characteristic_type ON action_logs(characteristic_type);
	`

	if _, err := al.db.Exec(query); err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

func (al *ActionLog) LogAction(entityID, characteristicType, newValue string) error {
	// Check if the new value is the same as the last logged value
	lastValueQuery := `
	SELECT new_value FROM action_logs
	WHERE entity_id = ? AND characteristic_type = ?
	ORDER BY timestamp DESC
	LIMIT 1
	`

	var lastValue string
	err := al.db.QueryRow(lastValueQuery, entityID, characteristicType).Scan(&lastValue)
	if err == nil && lastValue == newValue {
		slog.Debug("[ActionLog] Skipping duplicate value", "entityID", entityID, "type", characteristicType, "value", newValue)
		return nil
	}

	query := `
	INSERT INTO action_logs (entity_id, characteristic_type, new_value, timestamp)
	VALUES (?, ?, ?, ?)
	`

	_, err = al.db.Exec(query, entityID, characteristicType, newValue, time.Now())
	if err != nil {
		slog.Error("[ActionLog] Failed to log action", "error", err, "entityID", entityID, "type", characteristicType)
		return fmt.Errorf("failed to log action: %w", err)
	}

	slog.Debug("[ActionLog] Logged action", "entityID", entityID, "type", characteristicType, "newValue", newValue)
	return nil
}

func (al *ActionLog) LogAutomationEnabled(automationID int, enabled bool) error {
	entityID := fmt.Sprintf("%d", automationID)
	value := "disable"
	if enabled {
		value = "enable"
	}
	return al.LogAction(entityID, "automation", value)
}

func (al *ActionLog) GetAutomationEnabled(automationID int, defaultValue bool) bool {
	entityID := fmt.Sprintf("%d", automationID)

	query := `
	SELECT new_value FROM action_logs
	WHERE entity_id = ? AND characteristic_type = 'automation'
	ORDER BY timestamp DESC
	LIMIT 1
	`

	var value string
	err := al.db.QueryRow(query, entityID).Scan(&value)
	if err != nil {
		// No record found, return default
		return defaultValue
	}

	return value == "enable"
}

func (al *ActionLog) GetRecentLogs(limit int) ([]LogEntry, error) {
	query := `
	SELECT id, entity_id, characteristic_type, new_value, timestamp
	FROM action_logs
	ORDER BY timestamp DESC
	LIMIT ?
	`

	rows, err := al.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry

		err := rows.Scan(&log.ID, &log.EntityID, &log.CharacteristicType,
			&log.NewValue, &log.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (al *ActionLog) GetLogsByEntity(entityID string, limit int) ([]LogEntry, error) {
	query := `
	SELECT id, entity_id, characteristic_type, new_value, timestamp
	FROM action_logs
	WHERE entity_id = ?
	ORDER BY timestamp DESC
	LIMIT ?
	`

	rows, err := al.db.Query(query, entityID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry

		err := rows.Scan(&log.ID, &log.EntityID, &log.CharacteristicType,
			&log.NewValue, &log.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		logs = append(logs, log)
	}

	return logs, nil
}

func (al *ActionLog) Close() error {
	return al.db.Close()
}

