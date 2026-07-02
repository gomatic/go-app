package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

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
