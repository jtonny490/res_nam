package config

import "os"

type Config struct {
	DatabaseURL   string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBSSLMode     string
	JWTSecret     string
	KijaniAPIKey  string
	KijaniBaseURL string
	KijaniMode    string
	Port          string
	MigrationsDir string
	CORSOrigin    string
	Seed          bool
}

func Load() Config {
	// Defaults target a local Postgres; Docker Compose overrides DB_HOST with its
	// internal service name. When DATABASE_URL is set it takes precedence and is
	// used for managed Postgres (e.g. Render) which requires TLS.
	c := Config{
		DBHost:        "localhost",
		DBPort:        "5432",
		DBUser:        "postgres",
		DBName:        "res_nam",
		DBSSLMode:     "disable",
		Port:          "8080",
		KijaniMode:    "mock",
		MigrationsDir: "migrations",
		CORSOrigin:    "*",
	}
	c.DatabaseURL = os.Getenv("DATABASE_URL")
	if v := os.Getenv("DB_HOST"); v != "" {
		c.DBHost = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		c.DBPort = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		c.DBUser = v
	}
	c.DBPassword = os.Getenv("DB_PASSWORD")
	if v := os.Getenv("DB_NAME"); v != "" {
		c.DBName = v
	}
	if v := os.Getenv("DB_SSLMODE"); v != "" {
		c.DBSSLMode = v
	}
	c.JWTSecret = os.Getenv("JWT_SECRET")
	c.KijaniAPIKey = os.Getenv("KIJANI_API_KEY")
	c.KijaniBaseURL = os.Getenv("KIJANI_BASE_URL")
	if v := os.Getenv("KIJANI_MODE"); v != "" {
		c.KijaniMode = v
	}
	c.Seed = os.Getenv("SEED_DATA") == "true"
	if v := os.Getenv("PORT"); v != "" {
		c.Port = v
	}
	if v := os.Getenv("MIGRATIONS_DIR"); v != "" {
		c.MigrationsDir = v
	}
	if v := os.Getenv("CORS_ORIGIN"); v != "" {
		c.CORSOrigin = v
	}
	return c
}
