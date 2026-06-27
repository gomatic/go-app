# go-app

The urfave/cli/v3 framework shared by every SkyKernel CLI (package `app`): the `Runner`/`Default`/`action` combinators, the logger-in-metadata convention (`GetLogger`, `LoggerBefore`), the standard global flags (`LogLevelFlag`, `LogFormatFlag`, `OutputFlag`), and the `Run` harness. It composes `go-error` + `go-log` + `go-output`.

- **This is the only library that imports a CLI framework.** It holds nothing tied to any one command — command-specific flags (e.g. an endpoint flag), error values, and command trees stay in their own repos (`skytl`, `skykerneld`).
- Quality gate: gofumpt, `go vet`, staticcheck, govulncheck, gocognit ≤ 7, **100% coverage**.
