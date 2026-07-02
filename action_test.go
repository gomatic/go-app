package app

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/gomatic/go-output"
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
	require.New(t).ErrorIs(secondErr(runCommand(okRunner, "--output", "xml", "do")), output.ErrUnsupportedFormat)
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

func TestGetLoggerIndirection(t *testing.T) {
	t.Parallel()
	// The action path resolves its logger through the getLogger indirection.
	assert.New(t).NotNil(getLogger(nil))
}
