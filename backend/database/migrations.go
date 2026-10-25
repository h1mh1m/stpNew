package databases

import (
	"database/sql"
	"fmt"
	"log"
	"stpNew/backend/migrations"
)

// LoadMigrationSQL loads the authoritative migration SQL for a given version and direction.
// The SQL is retrieved directly from embedded migration files in backend/migrations/*.sql.
func LoadMigrationSQL(version int, direction string) (string, error) {
	var filename string
	switch version {
	case 1:
		if direction == "up" {
			filename = "000001_booking_integrity.up.sql"
		} else if direction == "down" {
			filename = "000001_booking_integrity.down.sql"
		} else {
			return "", fmt.Errorf("invalid migration direction: %s", direction)
		}
	default:
		return "", fmt.Errorf("unknown migration version: %d", version)
	}

	content, err := migrations.FS.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read embedded migration file %s: %w", filename, err)
	}
	return string(content), nil
}

// EnsureSchemaMigrationsTable creates the migrations tracker table if not exists.
func EnsureSchemaMigrationsTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS schema_migrations (
		version INT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := db.Exec(query)
	return err
}

// GetMigrationStatus returns the highest applied migration version (0 if none).
func GetMigrationStatus(db *sql.DB) (int, error) {
	if err := EnsureSchemaMigrationsTable(db); err != nil {
		return 0, fmt.Errorf("failed to ensure schema_migrations table: %w", err)
	}

	var version int
	err := db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&version)
	if err != nil {
		return 0, fmt.Errorf("failed to query migration version: %w", err)
	}
	return version, nil
}

// ApplyMigrationsUp executes migration version 1 inside a database transaction.
func ApplyMigrationsUp(db *sql.DB) error {
	currentVersion, err := GetMigrationStatus(db)
	if err != nil {
		return err
	}

	if currentVersion >= 1 {
		log.Println("Migrations already up to date (version 1)")
		return nil
	}

	sqlContent, err := LoadMigrationSQL(1, "up")
	if err != nil {
		return fmt.Errorf("failed to load migration up: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin migration transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(sqlContent); err != nil {
		return fmt.Errorf("migration up failed: %w", err)
	}

	if _, err := tx.Exec("INSERT INTO schema_migrations (version, applied_at) VALUES (1, CURRENT_TIMESTAMP)"); err != nil {
		return fmt.Errorf("failed to record migration version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration transaction: %w", err)
	}

	log.Println("Applied migration 000001_booking_integrity up successfully")
	return nil
}

// ApplyMigrationsDown rolls back migration version 1 inside a database transaction.
func ApplyMigrationsDown(db *sql.DB) error {
	currentVersion, err := GetMigrationStatus(db)
	if err != nil {
		return err
	}

	if currentVersion < 1 {
		log.Println("No migrations to roll back (version is 0)")
		return nil
	}

	sqlContent, err := LoadMigrationSQL(1, "down")
	if err != nil {
		return fmt.Errorf("failed to load migration down: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin rollback transaction: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.Exec(sqlContent); err != nil {
		return fmt.Errorf("migration down failed: %w", err)
	}

	if _, err := tx.Exec("DELETE FROM schema_migrations WHERE version = 1"); err != nil {
		return fmt.Errorf("failed to remove migration version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit rollback transaction: %w", err)
	}

	log.Println("Rolled back migration 000001_booking_integrity down successfully")
	return nil
}
