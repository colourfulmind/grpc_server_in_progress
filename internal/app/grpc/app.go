package grpcserver

import (
	"google.golang.org/grpc"
	"log/slog"
)

type App struct {
	log        *slog.Logger
	GRPCServer *grpc.Server
	host       string
	post       int
}

func New(log *slog.Logger, authService auth.Auth, blogService blog.Blog, host string, port int) *App {
	GRPCServer := grpc.NewServer()
	auth.Register(GRPCServer, authService)
	blog.Register(GRPCServer, blogService)
	return &App{
		log:        log,
		GRPCServer: GRPCServer,
		host:       host,
		post:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "internal.app.grpc.New()"
	return nil
}
