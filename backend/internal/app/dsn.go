package app

import (
	"fmt"
	"strings"

	"res_nam/internal/config"
)

// buildDSN returns the Postgres connection string. DATABASE_URL takes
// precedence (managed providers such as Render require TLS, so sslmode=require
// is added when absent). Otherwise the discrete DB_* settings are used.
func buildDSN(c config.Config) string {
	if c.DatabaseURL != "" {
		return withSSLMode(c.DatabaseURL, "require")
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

func withSSLMode(dsn, mode string) string {
	if strings.Contains(dsn, "sslmode=") {
		return dsn
	}
	separator := "?"
	if strings.Contains(dsn, "?") {
		separator = "&"
	}
	return dsn + separator + "sslmode=" + mode
}
