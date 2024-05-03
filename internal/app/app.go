package app

import (
	"log/slog"
	grpcserver "main/internal/app/grpc"
	"main/internal/config"
)

type App struct {
	Server *grpcserver.App
}

// New should return DB error
func New(log *slog.Logger, cfg *config.Config) *App {
	authService := grpcserver.New(log, cfg.Grpc.Host, cfg.Grpc.Port)
	return &App{
		//Server: server,
	}
}
