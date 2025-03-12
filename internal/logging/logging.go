package logging

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"slices"
	"time"
)

var nowFunc = time.Now

func Error(ctx context.Context, err error, attrs ...slog.Attr) {
	as := slices.Clone(attrs)
	as = append(as,
		slog.String("error.message", fmt.Sprintf("%s", err)),
		slog.String("error.type", fmt.Sprintf("%T", err)),
	)
	logAttrs(ctx, slog.Default(), 1, slog.LevelError, err.Error(), as...)
}

func logAttrs(ctx context.Context, logger *slog.Logger, nestLevel int, logLevel slog.Level, msg string, attrs ...slog.Attr) {
	if !logger.Enabled(ctx, logLevel) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2+nestLevel, pcs[:])
	r := slog.NewRecord(nowFunc(), logLevel, msg, pcs[0])
	r.AddAttrs(attrs...)
	_ = logger.Handler().Handle(ctx, r) //nolint:errcheck
}
