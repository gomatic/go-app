package pg

// Named types for each Config field, bound by the CLI via pointer conversion.
// They are unexported; the Config struct and its fields are exported because
// the CLI binds them.
type (
	host     string // host is the PostgreSQL server host (--pghost / PGHOST).
	port     int    // port is the PostgreSQL server port (--pgport / PGPORT).
	name     string // name is the database name (--pgdatabase / PGDATABASE).
	user     string // user is the database user (--pguser / PGUSER).
	password string // password is the database password (--pgpassword / PGPASSWORD).
	sslMode  string // sslMode is the libpq SSL mode (--pgsslmode / PGSSLMODE).
)
