package pg

import (
	"context"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestBinderFlags(t *testing.T) {
	var cfg Config
	got := Binder{Config: &cfg}.Flags()

	wantEnv := map[string]string{
		HostFlag:     "PGHOST",
		PortFlag:     "PGPORT",
		DatabaseFlag: "PGDATABASE",
		UserFlag:     "PGUSER",
		PasswordFlag: "PGPASSWORD",
		SSLModeFlag:  "PGSSLMODE",
	}
	if len(got) != len(wantEnv) {
		t.Fatalf("got %d flags, want %d", len(got), len(wantEnv))
	}
	seen := map[string]bool{}
	for _, fl := range got {
		flagName := fl.Names()[0]
		seen[flagName] = true
		env, ok := wantEnv[flagName]
		if !ok {
			t.Fatalf("unexpected flag %q", flagName)
		}
		keys := envKeys(t, fl)
		if len(keys) != 1 || keys[0] != env {
			t.Fatalf("flag %q env = %v, want [%s]", flagName, keys, env)
		}
	}
	for flagName := range wantEnv {
		if !seen[flagName] {
			t.Fatalf("missing flag %q", flagName)
		}
	}
}

// TestBinderDestinations proves parsed flag values land in the bound Config.
func TestBinderDestinations(t *testing.T) {
	var cfg Config
	cmd := &cli.Command{
		Name:   "test",
		Flags:  Binder{Config: &cfg}.Flags(),
		Action: func(context.Context, *cli.Command) error { return nil },
	}
	err := cmd.Run(context.Background(), []string{
		"test",
		"--pghost", "db.local",
		"--pgport", "5432",
		"--pgdatabase", "appdb",
		"--pguser", "app",
		"--pgpassword", "secret",
		"--pgsslmode", "verify-full",
	})
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	want := "host=db.local port=5432 user=app password=secret dbname=appdb sslmode=verify-full"
	if got := cfg.ConnString(); got != want {
		t.Fatalf("bound ConnString = %q, want %q", got, want)
	}
}

// envKeys extracts the env var keys from a string or int flag's Sources.
func envKeys(t *testing.T, fl cli.Flag) []string {
	t.Helper()
	switch f := fl.(type) {
	case *cli.StringFlag:
		return f.Sources.EnvKeys()
	case *cli.IntFlag:
		return f.Sources.EnvKeys()
	default:
		t.Fatalf("unexpected flag type %T", fl)
		return nil
	}
}
