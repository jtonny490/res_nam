package database

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

const trackingTable = "schema_migrations"

// Run applies every pending *.up.sql migration found in dir, in lexical order,
// recording applied versions in the schema_migrations table. Each migration is
// applied inside its own transaction so a failure leaves the version unrecorded.
func Run(db *gorm.DB, dir string) error {
	if err := db.Exec("CREATE TABLE IF NOT EXISTS " + trackingTable + " (version TEXT PRIMARY KEY, applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP)").Error; err != nil {
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	var ups []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			ups = append(ups, e.Name())
		}
	}
	sort.Strings(ups)
	for _, name := range ups {
		version := strings.TrimSuffix(name, ".up.sql")
		var count int64
		if err := db.Raw("SELECT COUNT(1) FROM "+trackingTable+" WHERE version = ?", version).Scan(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		body, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return err
		}
		if err := db.Transaction(func(tx *gorm.DB) error {
			for _, stmt := range statements(string(body)) {
				if err := tx.Exec(stmt).Error; err != nil {
					return err
				}
			}
			return tx.Exec("INSERT INTO "+trackingTable+" (version) VALUES (?)", version).Error
		}); err != nil {
			return err
		}
	}
	return nil
}

func statements(sql string) []string {
	parts := strings.Split(sql, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return out
}
