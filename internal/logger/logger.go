package logger

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fatih/color"
)

// Infof is an example of a user-defined logging function that wraps slog.
// The log record contains the source position of the caller of Infof.
func Logf(logger *slog.Logger, level slog.Level, format string, args ...any) {
	if !logger.Enabled(context.Background(), slog.LevelInfo) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]

	msg := fmt.Sprintf(format, args...)
	switch level {
	case slog.LevelDebug:
		msg = color.MagentaString(msg)
	case slog.LevelInfo:
		msg = color.BlueString(msg)
	case slog.LevelWarn:
		msg = color.YellowString(msg)
	case slog.LevelError:
		msg = color.RedString(msg)
	}

	r := slog.NewRecord(time.Now(), level, msg, pcs[0])
	_ = logger.Handler().Handle(context.Background(), r)
}

func New() *slog.Logger {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		// Remove time.
		if a.Key == slog.TimeKey && len(groups) == 0 {
			return slog.Attr{}
		}
		// Remove the directory from the source's filename.
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.File = filepath.Base(source.File)
		}
		return a
	}
	logger := slog.New(slog.NewTextHandler(color.Output, &slog.HandlerOptions{AddSource: false, ReplaceAttr: replace}))
	return logger
}
