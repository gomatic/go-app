package app

import (
	"context"
	"testing"

	"github.com/gomatic/go-log"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

func TestGlobalFlags(t *testing.T) {
	t.Parallel()
	want := assert.New(t)
	var cfg log.LoggerConfig
	flags := LoggerFlags{Config: &cfg, EnvPrefix: "APP_"}

	level := flags.LevelFlag().(*cli.StringFlag)
	want.Equal("log-level", level.Name)
	want.Equal("info", level.Value)
	want.Equal([]string{"APP_LOG_LEVEL"}, level.Sources.EnvKeys())

	format := flags.FormatFlag(log.FormatJSON).(*cli.StringFlag)
	want.Equal("log-format", format.Name)
	want.Equal("json", format.Value)
	want.Equal([]string{"APP_LOG_FORMAT"}, format.Sources.EnvKeys())

	out := OutputFlag("APP_").(*cli.StringFlag)
	want.Equal("output", out.Name)
	want.Equal([]string{"o"}, out.Aliases)
	want.Equal("json", out.Value)
}

// TestGlobalFlagsBind proves the level/format flags actually write through to the
// bound config when parsed from the command line.
func TestGlobalFlagsBind(t *testing.T) {
	t.Parallel()
	want := assert.New(t)
	var cfg log.LoggerConfig
	flags := LoggerFlags{Config: &cfg, EnvPrefix: "APP_"}
	root := &cli.Command{
		Name:     "root",
		Metadata: map[string]any{},
		Flags: []cli.Flag{
			flags.LevelFlag(),
			flags.FormatFlag(log.FormatText),
		},
		Action: func(context.Context, *cli.Command) error { return nil },
	}
	want.NoError(root.Run(context.Background(), []string{"root", "--log-level", "debug", "--log-format", "json"}))
	want.Equal(log.Level("debug"), cfg.Level)
	want.Equal(log.Format("json"), cfg.Format)
}
