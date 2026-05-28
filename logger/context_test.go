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
	adapter1 := NewSlogAdapter(slog.New(slog.NewJSONHandler(buf1, nil)))
	adapter2 := NewSlogAdapter(slog.New(slog.NewJSONHandler(buf2, nil)))

	ctx := NewContext(context.Background(), adapter1)
	ctx = NewContext(ctx, adapter2)

	extracted := FromContext(ctx)
	assert.Equal(t, adapter2, extracted)
}
