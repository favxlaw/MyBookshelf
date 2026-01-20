package store

import (
	"database/sql"
	"fmt"
	"log"
)

// Migration represents a database migration
type Migration struct {
	Version     int
	Description string
	Up          string
	Down        string
}

// migrations is the list of all migrations in order
var migrations = []Migration{
	{
		Version:     1,
		Description: "Create books table",
		Up: `
			CREATE TABLE IF NOT EXISTS books (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				title TEXT NOT NULL,
				author TEXT NOT NULL,
				status TEXT NOT NULL DEFAULT 'to_read',
				category TEXT,
				notes TEXT,
				start_date DATETIME NOT NULL,
				end_date DATETIME
			);
		`,
		Down: `DROP TABLE IF EXISTS books;`,
	},
	{
		Version:     2,
		Description: "Create schema_migrations table",
		Up: `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version INTEGER PRIMARY KEY,
				applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
			);
		`,
		Down: `DROP TABLE IF EXISTS schema_migrations;`,
	},
}

// RunMigrations executes all pending migrations
func RunMigrations(db *sql.DB) error {
	// Ensure schema_migrations table exists
	_, err := db.Exec(migrations[1].Up)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get current version
	currentVersion := getCurrentVersion(db)
	log.Printf("Current database version: %d", currentVersion)

	// Run pending migrations
	for _, migration := range migrations {
		if migration.Version <= currentVersion {
			continue // Already applied
		}

		log.Printf("Applying migration %d: %s", migration.Version, migration.Description)

		// Execute the migration
		_, err := db.Exec(migration.Up)
		if err != nil {
			return fmt.Errorf("migration %d failed: %w", migration.Version, err)
		}

		// Record that we applied this migration
		_, err = db.Exec(
			`INSERT INTO schema_migrations (version) VALUES (?)`,
			migration.Version,
		)
		if err != nil {
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		log.Printf("Migration %d applied successfully", migration.Version)
	}

	log.Println("All migrations up to date")
	return nil
}

// getCurrentVersion returns the latest applied migration version
func getCurrentVersion(db *sql.DB) int {
	var version int
	err := db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_migrations`).Scan(&version)
	if err != nil {
		return 0 // Table doesn't exist yet or is empty
	}
	return version
}

// RollbackMigration rolls back the last migration (optional, for learning)
func RollbackMigration(db *sql.DB) error {
	currentVersion := getCurrentVersion(db)
	if currentVersion == 0 {
		return fmt.Errorf("no migrations to rollback")
	}

	// Find the migration to rollback
	var targetMigration *Migration
	for i := range migrations {
		if migrations[i].Version == currentVersion {
			targetMigration = &migrations[i]
			break
		}
	}

	if targetMigration == nil {
		return fmt.Errorf("migration %d not found", currentVersion)
	}

	log.Printf("Rolling back migration %d: %s", targetMigration.Version, targetMigration.Description)

	// Execute the down migration
	_, err := db.Exec(targetMigration.Down)
	if err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	// Remove from schema_migrations
	_, err = db.Exec(`DELETE FROM schema_migrations WHERE version = ?`, currentVersion)
	if err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}

	log.Printf("Migration %d rolled back successfully", currentVersion)
	return nil
}
