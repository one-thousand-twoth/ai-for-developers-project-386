package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"ariga.io/atlas/sql/migrate"
)

//go:embed migrations/*.sql migrations/atlas.sum
var migrations embed.FS

// applyMigrations validates the committed Atlas directory and the applied prefix.
// DDL and revision records commit together under the driver's IMMEDIATE lock.
func applyMigrations(ctx context.Context, db *sql.DB) error {
	dir := &migrate.MemDir{}
	entries, err := fs.ReadDir(migrations, "migrations")
	if err != nil {
		return err
	}
	for _, entry := range entries {
		body, err := migrations.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return err
		}
		if err := dir.WriteFile(entry.Name(), body); err != nil {
			return err
		}
	}
	if err := migrate.Validate(dir); err != nil {
		return fmt.Errorf("migration checksum: %w", err)
	}
	files, err := dir.Files()
	if err != nil {
		return err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version TEXT PRIMARY KEY, hash TEXT NOT NULL)`); err != nil {
		return err
	}
	rows, err := tx.QueryContext(ctx, `SELECT version, hash FROM schema_migrations ORDER BY version`)
	if err != nil {
		return err
	}
	applied := 0
	for rows.Next() {
		var version, hash string
		if err := rows.Scan(&version, &hash); err != nil {
			_ = rows.Close()
			return err
		}
		if applied >= len(files) || version != files[applied].Version() || hash != migrationHash(files[applied]) {
			_ = rows.Close()
			return fmt.Errorf("migration history differs at version %s", version)
		}
		applied++
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	for _, file := range files[applied:] {
		statements, err := file.Stmts()
		if err != nil {
			return fmt.Errorf("parse %s: %w", file.Name(), err)
		}
		for _, statement := range statements {
			if _, err := tx.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("apply %s: %w", file.Name(), err)
			}
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version, hash) VALUES (?, ?)`, file.Version(), migrationHash(file)); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func migrationHash(file migrate.File) string {
	return fmt.Sprintf("%x", sha256.Sum256(file.Bytes()))
}
