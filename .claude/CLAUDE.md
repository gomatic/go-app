# go-app

The urfave/cli/v3 framework for CLIs built to gomatic/template.cli standards (package `app`): the `Runner`/`Default`/`action` combinators, logger-in-metadata (`GetLogger`, `LoggerBefore`), standard global flags (`LoggerFlags` with `LevelFlag`/`FormatFlag`, `OutputFlag`), and the `Run` harness. Composes `go-log` + `go-output` + `urfave/cli/v3`. Generic — lives in `gomatic`, consumed by `template.cli` and the gomatic CLIs.

- `flags/pg` (package `pg`): the standard libpq `PG*` flag set (`pg.Binder`, same binder pattern as `LoggerFlags`) and the `pg.Config` it binds, with `ConnString()` emitting a libpq KV string whose empty fields defer to the `PG*` env vars; opening a pool stays with the consumer. Extracted from xto-email's repo split (see `xto-email/_projects/specs/repo-split/`).
- The only lib importing a CLI framework. Holds nothing tied to any one command.
- Gate: gofumpt, vet, staticcheck, govulncheck, gocognit ≤ 7, 100% coverage.
