package database

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestRunAppliesTracksAndIsIdempotent(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("000001_create_widgets.up.sql", "CREATE TABLE widgets (id INTEGER PRIMARY KEY, name TEXT); CREATE INDEX widgets_name_idx ON widgets(name);")
	write("000002_seed_widget.up.sql", "INSERT INTO widgets (name) VALUES ('a');")

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)

	if err := Run(db, dir); err != nil {
		t.Fatalf("first run: %v", err)
	}
	if err := Run(db, dir); err != nil {
		t.Fatalf("second run: %v", err)
	}

	var versions, widgets int64
	db.Raw("SELECT COUNT(1) FROM schema_migrations").Scan(&versions)
	db.Raw("SELECT COUNT(1) FROM widgets").Scan(&widgets)
	if versions != 2 {
		t.Fatalf("tracked migrations = %d, want 2", versions)
	}
	if widgets != 1 {
		t.Fatalf("widgets = %d, want 1 (migrations must not re-apply)", widgets)
	}
}
