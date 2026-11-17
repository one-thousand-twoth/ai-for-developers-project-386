package storage

import (
	"bytes"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"ariga.io/atlas/sql/migrate"
)

func TestMigrationsChecksumAfterCRLFCheckout(t *testing.T) {
	entries, err := fs.ReadDir(migrations, "migrations")
	if err != nil {
		t.Fatal(err)
	}
	files := make(map[string][]byte, len(entries))
	for _, entry := range entries {
		body, err := migrations.ReadFile("migrations/" + entry.Name())
		if err != nil {
			t.Fatal(err)
		}
		files["internal/storage/migrations/"+entry.Name()] = body
	}
	directory := checkoutWithAutoCRLF(t, files)
	dir, err := migrate.NewLocalDir(filepath.Join(directory, "internal", "storage", "migrations"))
	if err != nil {
		t.Fatal(err)
	}
	if err := migrate.Validate(dir); err != nil {
		t.Fatalf("migration checksum after core.autocrlf checkout: %v", err)
	}
}

func TestCheckoutPreservesProjectLineEndings(t *testing.T) {
	text := []byte("first line\nsecond line\n")
	files := map[string][]byte{
		".editorconfig":                 text,
		"README.md":                     text,
		"Dockerfile":                    text,
		".github/workflows/codegen.yml": text,
		"cmd/server/main.go":            text,
		"gen/api/oas_schemas_gen.go":    text,
		"web/src/App.tsx":               text,
		"web/public/favicon.svg":        text,
		"web/src/assets/hero.png":       {'P', 'N', 'G', 0, '\r', '\n'},
	}
	directory := checkoutWithAutoCRLF(t, files)
	for name, want := range files {
		t.Run(name, func(t *testing.T) {
			got, err := os.ReadFile(filepath.Join(directory, name))
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("checkout changed file bytes: got %q, want %q", got, want)
			}
		})
	}
}

func checkoutWithAutoCRLF(t *testing.T, files map[string][]byte) string {
	t.Helper()
	git, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git is required to test checkout line endings")
	}
	attributes, err := os.ReadFile(filepath.Join("..", "..", ".gitattributes"))
	if err != nil {
		t.Fatal(err)
	}
	directory := t.TempDir()
	if err := os.WriteFile(filepath.Join(directory, ".gitattributes"), attributes, 0o644); err != nil {
		t.Fatal(err)
	}
	for name, body := range files {
		path := filepath.Join(directory, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, body, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	runGit := func(args ...string) {
		t.Helper()
		args = append([]string{"-c", "core.autocrlf=true", "-c", "core.safecrlf=false"}, args...)
		command := exec.Command(git, args...)
		command.Dir = directory
		command.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+filepath.Join(directory, "gitconfig"))
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, output)
		}
	}
	runGit("init", "--quiet")
	runGit("add", "--", ".")
	for name := range files {
		if err := os.Remove(filepath.Join(directory, name)); err != nil {
			t.Fatal(err)
		}
	}
	runGit("checkout-index", "--all")
	return directory
}
