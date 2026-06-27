# go-app

The urfave/cli/v3 framework for CLIs built to gomatic/template.cli standards (package `app`): the `Runner`/`Default`/`action` combinators, logger-in-metadata (`GetLogger`, `LoggerBefore`), standard global flags (`LogLevelFlag`, `LogFormatFlag`, `OutputFlag`), and the `Run` harness. Composes `go-error` + `go-log` + `go-output`. Generic — lives in `gomatic`, consumed by `template.cli` and the SkyKernel CLIs.

- The only lib importing a CLI framework. Holds nothing tied to any one command.
- Gate: gofumpt, vet, staticcheck, govulncheck, gocognit ≤ 7, 100% coverage.
