package logger

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext_NestedOverwrite(t *testing.T) {
	buf1 := &bytes.Buffer{}
	buf2 := &bytes.Buffer{}
	logger1 := slog.New(slog.NewJSONHandler(buf1, nil))
	logger2 := slog.New(slog.NewJSONHandler(buf2, nil))

	ctx := NewContext(context.Background(), logger1)
	ctx = NewContext(ctx, logger2)

	extracted := FromContext(ctx)
	assert.Equal(t, logger2, extracted)
}
