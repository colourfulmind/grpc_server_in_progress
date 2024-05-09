package app

import (
	"fmt"
	"log/slog"
	grpcserver "main/internal/app/grpc"
	"main/internal/config"
	"main/internal/services/sso"
	"main/internal/storage/postgres"
	"main/pkg/logger/sl"
)

// App is a structure to store a pointer to grpc server
type App struct {
	Server *grpcserver.App
}

// New connects to postgres DB, register services and returns a pointer to created server.
// Returns error if connection to DB failed.
func New(log *slog.Logger, cfg *config.Config) (*App, error) {
	const op = "internal.app.New"

	storage, err := postgres.New(cfg.Postgres)
	if err != nil {
		log.Error("error occurred while connecting to database", slog.String("op", op), sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ssoService := sso.New(log, storage, storage, cfg.TokenTTl)

	return &App{
		Server: grpcserver.New(ssoService, log, cfg.Grpc.Host, cfg.Grpc.Port),
	}, nil
}
