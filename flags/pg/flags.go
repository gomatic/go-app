package pg

import "github.com/urfave/cli/v3"

// PG* flag names (constants for reuse in tests/help).
const (
	HostFlag     = "pghost"
	PortFlag     = "pgport"
	DatabaseFlag = "pgdatabase"
	UserFlag     = "pguser"
	PasswordFlag = "pgpassword"
	SSLModeFlag  = "pgsslmode"
)

// Binder binds the standard libpq-compatible PG* flags to a shared Config.
// The destination pointer lives in a field (not a parameter) so the binder
// itself travels by value while the flags still write into the one Config
// instance the command's action reads — the same binder pattern as
// app.LoggerFlags.
type Binder struct {
	Config *Config
}

// Flags builds the PG* flag set bound to the Config destination. The env vars
// are the canonical PG* names (not app-prefixed) so the flags work with psql,
// pgx, and any libpq tool.
func (b Binder) Flags() []cli.Flag {
	cfg := b.Config
	return []cli.Flag{
		&cli.StringFlag{
			Name:        HostFlag,
			Sources:     cli.EnvVars("PGHOST"),
			Usage:       "PostgreSQL server host",
			Destination: (*string)(&cfg.Host),
		},
		&cli.IntFlag{
			Name:        PortFlag,
			Sources:     cli.EnvVars("PGPORT"),
			Usage:       "PostgreSQL server port",
			Destination: (*int)(&cfg.Port),
		},
		&cli.StringFlag{
			Name:        DatabaseFlag,
			Sources:     cli.EnvVars("PGDATABASE"),
			Usage:       "PostgreSQL database name",
			Destination: (*string)(&cfg.Name),
		},
		&cli.StringFlag{
			Name:        UserFlag,
			Sources:     cli.EnvVars("PGUSER"),
			Usage:       "PostgreSQL user name",
			Destination: (*string)(&cfg.User),
		},
		&cli.StringFlag{
			Name:        PasswordFlag,
			Sources:     cli.EnvVars("PGPASSWORD"),
			Usage:       "PostgreSQL password (optional — prefer ~/.pgpass)",
			Destination: (*string)(&cfg.Password),
		},
		&cli.StringFlag{
			Name:        SSLModeFlag,
			Sources:     cli.EnvVars("PGSSLMODE"),
			Usage:       "PostgreSQL SSL mode (disable, require, verify-ca, verify-full)",
			Destination: (*string)(&cfg.SSLMode),
		},
	}
}
