package app

import (
	"log/slog"
	"os"
)

var DebugMode = os.Getenv("VISION_DEBUG") == "1"

// Debug logs debug messages only when VISION_DEBUG=1 is set.
// Use this for hot-path logging that could impact performance.
func Debug(msg string, args ...interface{}) {
	if DebugMode {
		slog.Debug(msg, args...)
	}
}

// Info logs informational messages unconditionally.
func Info(msg string, args ...interface{}) {
	slog.Info(msg, args...)
}

// Warn logs warning messages unconditionally.
func Warn(msg string, args ...interface{}) {
	slog.Warn(msg, args...)
}

// Error logs error messages unconditionally.
func Error(msg string, args ...interface{}) {
	slog.Error(msg, args...)
}
