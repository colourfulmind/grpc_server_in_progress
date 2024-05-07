package grpcserver

import (
	"main/internal/grpc/sso"
)

type App struct {
	GRPCServer *grpc.Server
	host       string
	port       int
}

func New(ssoServer sso.SSO, host string, port int) *App {
	GRPCServer := grpc.NewServer()
	sso.Register(GRPCServer, ssoServer, log)
	return &App{
		GRPCServer: GRPCServer,
		host:       host,
		port:       port,
	}
}
