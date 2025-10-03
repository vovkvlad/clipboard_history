package storage

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Migration struct {
	version int
	name    string
	sql     string
}

func ensureMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY, 
		name TEXT NOT NULL, 
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP)
	`)
	return err
}

func getMigrations() ([]Migration, error) {
	var migrations []Migration

	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		namePart := strings.Split(entry.Name(), ".")[0]
		// Extract version number from filename (e.g., "001_init.sql" -> 1)
		versionStr := strings.Split(namePart, "_")[0]
		version, err := strconv.Atoi(versionStr)
		migrationName := strings.Join(strings.Split(namePart, "_")[1:], "_")
		if err != nil {
			log.Printf("Invalid migration filename format: %s", entry.Name())
			return nil, err
		}

		content, err := fs.ReadFile(migrationsFS, filepath.Join("migrations", entry.Name()))
		if err != nil {
			log.Printf("Failed to read migration file %s: %v", entry.Name(), err)
			return nil, err
		}

		migration := Migration{
			version: version,
			name:    migrationName,
			sql:     string(content),
		}

		migrations = append(migrations, migration)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].version < migrations[j].version
	})

	return migrations, nil
}

func RunMigrations(db *sql.DB) error {
	if err := ensureMigrationsTable(db); err != nil {
		return err
	}

	migrations, err := getMigrations()
	if err != nil {
		log.Printf("Failed to get migrations: %v", err)
		return err
	}
	log.Printf("Found %d migrations", len(migrations))

	// Get already applied migrations
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		log.Printf("Failed to get applied migrations: %v", err)
		return err
	}

	for _, migration := range migrations {
		if !appliedMigrations[migration.version] {
			log.Printf("Applying migration: %s", migration.name)

			if _, err := db.Exec(migration.sql); err != nil {
				log.Printf("Failed to apply migration %s: %v", migration.name, err)
				return err
			}

			if err := recordMigrationApplied(db, migration.name); err != nil {
				log.Printf("Failed to record migration %s: %v", migration.name, err)
				return err
			}
		}
	}

	return nil
}

func getAppliedMigrations(db *sql.DB) (map[int]bool, error) {
	rows, err := db.Query("SELECT id FROM migrations ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}()

	applied := make(map[int]bool)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		applied[id] = true
	}

	return applied, rows.Err()
}

func recordMigrationApplied(db *sql.DB, name string) error {
	_, err := db.Exec("INSERT INTO migrations (name) VALUES (?)", name)
	return err
}
