package logging

import "log/slog"

func New(_ string, _ string) *slog.Logger {
	return slog.Default()
}
