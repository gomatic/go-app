// Package pg supplies the standard libpq PG* flag set and the connection
// Config it binds. Consumers embed Config in a command config, list
// Binder{Config: &cfg}.Flags() in the command's flags, and pass
// Config.ConnString() to their PostgreSQL driver; empty fields are omitted so
// the driver falls back to the PG* environment variables (and ~/.pgpass /
// ~/.pg_service.conf). Opening a pool stays with the consumer.
package pg

import (
	"fmt"
	"strings"
)

// Config holds libpq-compatible PostgreSQL connection parameters, bound from
// the standard PG* flags. Fields are exported because the CLI binds them by
// pointer; the field types are unexported domain types (see types.go).
type Config struct {
	Host     host
	Name     name
	User     user
	Password password
	SSLMode  sslMode
	Port     port
}

// ConnString constructs a libpq KV connection string from the explicit flag
// values. Empty fields are omitted so the driver falls back to the PG*
// environment variables (and ~/.pgpass / ~/.pg_service.conf).
func (c Config) ConnString() string {
	var parts []string
	if c.Host != "" {
		parts = append(parts, fmt.Sprintf("host=%s", c.Host))
	}
	if c.Port != 0 {
		parts = append(parts, fmt.Sprintf("port=%d", c.Port))
	}
	if c.User != "" {
		parts = append(parts, fmt.Sprintf("user=%s", c.User))
	}
	if c.Password != "" {
		parts = append(parts, fmt.Sprintf("password=%s", c.Password))
	}
	if c.Name != "" {
		parts = append(parts, fmt.Sprintf("dbname=%s", c.Name))
	}
	if c.SSLMode != "" {
		parts = append(parts, fmt.Sprintf("sslmode=%s", c.SSLMode))
	}
	return strings.Join(parts, " ")
}
