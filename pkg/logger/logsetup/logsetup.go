package logsetup

import (
	"articles/pkg/logger/handlers/slogpretty"
	"log/slog"
	"os"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func SetupLogger(env string) *slog.Logger {
	var level slog.Leveler

	switch env {
	case EnvLocal:
		level = slog.LevelDebug
	case EnvDev:
		level = slog.LevelDebug
	case EnvProd:
		level = slog.LevelInfo
	}

	return SetupPrettySlog(level)
}

func SetupPrettySlog(level slog.Leveler) *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
