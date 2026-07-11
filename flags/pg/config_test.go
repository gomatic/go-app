package pg

import "testing"

func TestConnStringEmpty(t *testing.T) {
	if got := (Config{}).ConnString(); got != "" {
		t.Fatalf("empty Config ConnString = %q, want empty", got)
	}
}

func TestConnStringFull(t *testing.T) {
	cfg := Config{
		Host:     "db.local",
		Port:     26257,
		Name:     "app",
		User:     "app",
		Password: "secret",
		SSLMode:  "require",
	}
	want := "host=db.local port=26257 user=app password=secret dbname=app sslmode=require"
	if got := cfg.ConnString(); got != want {
		t.Fatalf("ConnString = %q, want %q", got, want)
	}
}
