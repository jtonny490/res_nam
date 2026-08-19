package config

import "os"

type Config struct{ DBHost, DBPort, DBUser, DBPassword, DBName, JWTSecret, KijaniAPIKey, Port string }

func Load() Config {
	c := Config{DBHost: "db", DBPort: "5432", DBUser: "postgres", DBName: "res_nam", Port: "8080"}
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
	c.JWTSecret = os.Getenv("JWT_SECRET")
	c.KijaniAPIKey = os.Getenv("KIJANI_API_KEY")
	if v := os.Getenv("PORT"); v != "" {
		c.Port = v
	}
	return c
}
