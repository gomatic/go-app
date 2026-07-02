package app

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v3"
)

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
