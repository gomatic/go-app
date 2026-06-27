# go-app

A [urfave/cli/v3](https://github.com/urfave/cli) application framework and run harness shared by gomatic CLIs.

`go-app` supplies the generic, command-agnostic glue — the `Runner`/`Default` action combinators, the logger-in-metadata convention, the standard global flags, and the signal-aware `Run` harness — and nothing tied to any one command. Command-specific flags, errors, and command trees stay in their own repositories. It composes [go-log](https://github.com/gomatic/go-log) and [go-output](https://github.com/gomatic/go-output) on top of urfave/cli/v3.

## Install

```sh
go get github.com/gomatic/go-app
```

## Usage

Define a `Runner` for your command's work, bind it into a cli action with `Default`, and run the root command through `Run` for signal-aware (SIGINT/SIGTERM) cancellation:

```go
package main

import (
	"context"
	"log/slog"
	"os"

	app "github.com/gomatic/go-app"
	"github.com/gomatic/go-log"
	"github.com/urfave/cli/v3"
)

type config struct{}

type result struct {
	Message string `json:"message" yaml:"message"`
}

// greet is the command's work: it receives the bound config and positional
// arguments and returns a result the action combinator encodes.
func greet(_ context.Context, logger *slog.Logger, _ config, args ...string) (result, error) {
	logger.Info("greeting", "args", args)
	return result{Message: "hello"}, nil
}

func main() {
	var cfg config
	var logCfg log.LoggerConfig

	cmd := &cli.Command{
		Name:     "greeter",
		Metadata: map[string]any{},
		Flags: []cli.Flag{
			app.LogLevelFlag(&logCfg, "GREETER_"),
			app.LogFormatFlag(&logCfg, "GREETER_", log.FormatText),
			app.OutputFlag("GREETER_"),
		},
		Before: app.LoggerBefore(func(*cli.Command) *slog.Logger {
			return logCfg.NewLogger(os.Stderr)
		}),
		Action: app.Default(&cfg, greet),
	}

	app.Run(context.Background(), cmd, os.Args, os.Exit)
}
```

The global flags resolve from `--flag`, the matching `GREETER_*` environment variables, or their defaults. `OutputFlag` selects the result encoding (`json` or `yaml`), which the `Default` action applies when writing the runner's result.

## API

- `Run(ctx, cmd, args, exit)` — runs a `*cli.Command` with SIGINT/SIGTERM cancellation, logging any error and exiting non-zero via the injected `exit`.
- `Default(cfg, runner)` — binds a config pointer and a `Runner` into a cli action that encodes the result via the output flag.
- `Runner[CONFIG, RESULT]` — the function type for a command's work.
- `LogLevelFlag`, `LogFormatFlag`, `OutputFlag` — the standard global flags.
- `GetLogger`, `LoggerBefore`, `LoggerMetadataKey` — the logger-in-metadata convention.

See the full reference with `go doc -all github.com/gomatic/go-app`.
