package sso

import (
	"context"
	server "main/protos/gen/go/blog"
)

type ServerSSO struct {
	server.UnimplementedSSOServer
	sso SSO
}

type SSO interface {
	RegisterNewUser(ctx context.Context, email, password string) (int64, error)
	Login(ctx context.Context, email, password string, appID int32) (string, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

func RegisterNewUser(ctx context.Context, req *server.RegisterRequest) (*server.RegisterResponse, error) {
	panic("not implemented")
}

func Login(ctx context.Context, req *server.LoginRequest) (*server.LoginResponse, error) {
	panic("not implemented")
}

func IsAdmin(ctx context.Context, req *server.AdminRequest) (*server.AdminResponse, error) {
	panic("not implemented")
}
