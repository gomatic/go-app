package app

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	slogx "github.com/skykernel/go-log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

type actionResult struct {
	V string `json:"v" yaml:"v"`
}

// runCommand builds a root command with the standard output flag and a single
// child whose action is built from Default, then runs it with the supplied args.
func runCommand(runner Runner[struct{}, actionResult], args ...string) (string, error) {
	var cfg struct{}
	var buf bytes.Buffer
	root := &cli.Command{
		Name:     "root",
		Writer:   &buf,
		Metadata: map[string]any{},
		Flags:    []cli.Flag{OutputFlag("ROOT_")},
		Commands: []*cli.Command{
			{Name: "do", Action: Default(&cfg, runner)},
		},
	}
	err := root.Run(context.Background(), append([]string{"root"}, args...))
	return buf.String(), err
}

func okRunner(_ context.Context, _ *slog.Logger, _ struct{}, _ ...string) (actionResult, error) {
	return actionResult{V: "ok"}, nil
}

func TestActionEncodesJSON(t *testing.T) {
	want, must := assert.New(t), require.New(t)
	out, err := runCommand(okRunner, "do")
	must.NoError(err)
	want.Equal("{\n  \"v\": \"ok\"\n}\n", out)
}

func TestActionEncodesYAML(t *testing.T) {
	want, must := assert.New(t), require.New(t)
	out, err := runCommand(okRunner, "--output", "yaml", "do")
	must.NoError(err)
	want.Equal("v: ok\n", out)
}

func TestActionUnsupportedFormat(t *testing.T) {
	require.New(t).ErrorContains(secondErr(runCommand(okRunner, "--output", "xml", "do")), "unsupported output format")
}

func TestActionPropagatesRunnerError(t *testing.T) {
	sentinel := errors.New("boom")
	failing := func(_ context.Context, _ *slog.Logger, _ struct{}, _ ...string) (actionResult, error) {
		return actionResult{}, sentinel
	}
	assert.New(t).ErrorIs(secondErr(runCommand(failing, "do")), sentinel)
}

// secondErr discards the output string and returns only the error, keeping the
// unsupported-format and runner-error cases to a single readable line.
func secondErr(_ string, err error) error { return err }

func TestGetLogger(t *testing.T) {
	t.Parallel()
	want := assert.New(t)

	want.NotNil(GetLogger(nil))

	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	withLogger := &cli.Command{Metadata: map[string]any{LoggerMetadataKey: logger}}
	want.Same(logger, GetLogger(withLogger))

	withoutLogger := &cli.Command{Metadata: map[string]any{}}
	want.NotNil(GetLogger(withoutLogger))
}

func TestLoggerBeforeStoresLogger(t *testing.T) {
	t.Parallel()
	want := assert.New(t)
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	cmd := &cli.Command{Metadata: map[string]any{}}
	before := LoggerBefore(func(*cli.Command) *slog.Logger { return logger })
	_, err := before(context.Background(), cmd)
	want.NoError(err)
	want.Same(logger, GetLogger(cmd))
}

func TestGetLoggerIndirection(t *testing.T) {
	t.Parallel()
	// The action path resolves its logger through the getLogger indirection.
	assert.New(t).NotNil(getLogger(nil))
}

func TestGlobalFlags(t *testing.T) {
	t.Parallel()
	want := assert.New(t)
	var cfg slogx.LoggerConfig

	level := LogLevelFlag(&cfg, "APP_").(*cli.StringFlag)
	want.Equal("log-level", level.Name)
	want.Equal("info", level.Value)
	want.Equal([]string{"APP_LOG_LEVEL"}, level.Sources.EnvKeys())

	format := LogFormatFlag(&cfg, "APP_", slogx.FormatJSON).(*cli.StringFlag)
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
	var cfg slogx.LoggerConfig
	root := &cli.Command{
		Name:     "root",
		Metadata: map[string]any{},
		Flags: []cli.Flag{
			LogLevelFlag(&cfg, "APP_"),
			LogFormatFlag(&cfg, "APP_", slogx.FormatText),
		},
		Action: func(context.Context, *cli.Command) error { return nil },
	}
	want.NoError(root.Run(context.Background(), []string{"root", "--log-level", "debug", "--log-format", "json"}))
	want.Equal(slogx.LogLevel("debug"), cfg.LogLevel)
	want.Equal(slogx.LogFormat("json"), cfg.LogFormat)
}

func TestRunSuccess(t *testing.T) {
	want := assert.New(t)
	exited := false
	cmd := &cli.Command{Name: "tool", Action: func(context.Context, *cli.Command) error { return nil }}
	Run(context.Background(), cmd, []string{"tool"}, func(int) { exited = true })
	want.False(exited, "successful run must not exit non-zero")
}

func TestRunError(t *testing.T) {
	want := assert.New(t)
	var code int
	cmd := &cli.Command{
		Name:   "tool",
		Action: func(context.Context, *cli.Command) error { return errors.New("boom") },
	}
	Run(context.Background(), cmd, []string{"tool"}, func(c int) { code = c })
	want.Equal(1, code)
}
