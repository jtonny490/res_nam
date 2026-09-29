package app

import (
	"testing"

	"res_nam/internal/config"
)

func TestBuildDSNPrefersDatabaseURLWithSSL(t *testing.T) {
	got := buildDSN(config.Config{DatabaseURL: "postgresql://u:p@host:5432/db"})
	want := "postgresql://u:p@host:5432/db?sslmode=require"
	if got != want {
		t.Fatalf("dsn = %q, want %q", got, want)
	}
}

func TestBuildDSNKeepsExistingSSLMode(t *testing.T) {
	dsn := "postgres://u:p@host/db?sslmode=disable"
	if got := buildDSN(config.Config{DatabaseURL: dsn}); got != dsn {
		t.Fatalf("dsn = %q, want %q", got, dsn)
	}
}

func TestBuildDSNAppendsWithAmpersand(t *testing.T) {
	got := buildDSN(config.Config{DatabaseURL: "postgres://u:p@host/db?application_name=x"})
	want := "postgres://u:p@host/db?application_name=x&sslmode=require"
	if got != want {
		t.Fatalf("dsn = %q, want %q", got, want)
	}
}

func TestBuildDSNUsesDiscreteFields(t *testing.T) {
	c := config.Config{DBHost: "db", DBPort: "5432", DBUser: "postgres", DBPassword: "pw", DBName: "res_nam", DBSSLMode: "disable"}
	want := "host=db port=5432 user=postgres password=pw dbname=res_nam sslmode=disable"
	if got := buildDSN(c); got != want {
		t.Fatalf("dsn = %q, want %q", got, want)
	}
}
