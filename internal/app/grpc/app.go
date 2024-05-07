package grpcserver

import (
	"main/internal/grpc/sso"
)

// App is a structure to store a host and a port to connect to server and a pointer to server itself.
type App struct {
	GRPCServer *grpc.Server
	host       string
	port       int
}

// New returns  connection to grpc server.
func New(ssoServer sso.SSO, host string, port int) *App {
	GRPCServer := grpc.NewServer()
	sso.Register(GRPCServer, ssoServer, log)
	return &App{
		GRPCServer: GRPCServer,
		host:       host,
		port:       port,
	}
}
