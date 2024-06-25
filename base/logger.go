package base

import (
	"log/slog"
	"os"
)

var LOG = slog.New(slog.NewJSONHandler(os.Stdout,
	&slog.HandlerOptions{Level: slog.LevelInfo}))
