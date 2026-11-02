// storage-diff creates migration files; it is deliberately not part of generate.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/one-thousand-twoth/calendar/internal/storage/ent/migrate"
	_ "modernc.org/sqlite"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	name := os.Getenv("NAME")
	if !regexp.MustCompile(`^[a-z][a-z0-9_]*$`).MatchString(name) {
		return fmt.Errorf("set NAME to a migration name, e.g. NAME=create_meeting_types mise run storage:diff")
	}
	dir, err := atlas.NewLocalDir("internal/storage/migrations")
	if err != nil {
		return err
	}
	// A private, single-connection database is only used to replay migration history.
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()
	db.SetMaxOpenConns(1)
	m, err := schema.NewMigrate(entsql.OpenDB(dialect.SQLite, db),
		schema.WithDir(dir), schema.WithMigrationMode(schema.ModeReplay), schema.WithFormatter(atlas.DefaultFormatter))
	if err != nil {
		return err
	}
	return m.NamedDiff(context.Background(), name, migrate.Tables...)
}
