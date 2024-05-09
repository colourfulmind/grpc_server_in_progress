package grpcserver

import (
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"main/internal/grpc/sso"
	"net"
)

// App is a structure to store a host and a port to connect to server and a pointer to server itself.
type App struct {
	GRPCServer *grpc.Server
	log        *slog.Logger
	host       string
	port       int
}

// New returns  connection to grpc server.
func New(ssoServer sso.SSO, log *slog.Logger, host string, port int) *App {
	GRPCServer := grpc.NewServer()
	sso.Register(GRPCServer, ssoServer, log)
	return &App{
		GRPCServer: GRPCServer,
		log:        log,
		host:       host,
		port:       port,
	}
}

// MustRun panics if grpc server wasn't started.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs grpc server. Returns error if it cannot serve the address.
func (a *App) Run() error {
	const op = "internal.app.grpc.Run"

	log := a.log.With(slog.String("op", op))
	log.Info("starting grpc server", slog.String("host", a.host), slog.Int("port", a.port))

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.host, a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err = a.GRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// Stop gracefully stops the application.
func (a *App) Stop() {
	const op = "internal.app.grpc.Stop"

	log := a.log.With(slog.String("op", op))
	log.Info("stopping grpc server", slog.String("host", a.host), slog.Int("port", a.port))

	a.GRPCServer.GracefulStop()

	log.Info("grpc server stopped")
}
