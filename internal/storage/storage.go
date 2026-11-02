// Package storage owns the Ent-backed SQLite database and versioned migrations.
package storage

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/one-thousand-twoth/calendar/internal/storage/ent"
	_ "modernc.org/sqlite"
)

// Store owns its connection pool. Every transaction starts with BEGIN IMMEDIATE.
type Store struct {
	client *ent.Client
}

func Open(ctx context.Context, path string) (*Store, error) {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(absolute), 0o700); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}
	u := url.URL{Scheme: "file", Path: filepath.ToSlash(absolute)}
	query := url.Values{
		"_pragma": {"foreign_keys(1)", "busy_timeout(5000)"},
		"_txlock": {"immediate"},
	}
	u.RawQuery = query.Encode()
	db, err := sql.Open("sqlite", u.String())
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	// A small bounded pool; correctness also holds across independent processes.
	db.SetMaxOpenConns(4)
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := applyMigrations(ctx, db); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("migrate database: %w", err)
	}
	return &Store{client: ent.NewClient(ent.Driver(entsql.OpenDB(dialect.SQLite, db)))}, nil
}

func (s *Store) Close() error { return s.client.Close() }
