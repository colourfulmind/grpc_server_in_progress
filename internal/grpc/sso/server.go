package sso

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"main/internal/ewrap"
	"main/internal/services/sso"
	"main/pkg/logger/sl"
	server "main/protos/gen/go/blog"
)

type ServerSSO struct {
	server.UnimplementedSSOServer
	sso SSO
	log *slog.Logger
}

type SSO interface {
	RegisterNewUser(ctx context.Context, email, password string) (int64, error)
	Login(ctx context.Context, email, password string, appID int32) (string, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

func Register(s *grpc.Server, sso SSO, log *slog.Logger) {
	server.RegisterSSOServer(s, &ServerSSO{
		sso: sso,
		log: log,
	})
}

func (s *ServerSSO) RegisterNewUser(ctx context.Context, req *server.RegisterRequest) (*server.RegisterResponse, error) {
	const op = "internal.grpc.sso.RegisterNewUser"

	log := s.log.With(slog.String("op", op), slog.String("email", req.GetEmail()))
	log.Info("register new user")

	if err := ValidateCredentials(req); err != nil {
		log.Error("invalid email or password", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := s.sso.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, sso.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))
			return nil, ewrap.UserAlreadyExists
		}
		if errors.Is(err, sso.ErrConnectionTime) {
			log.Error("failed to connect to the database", sl.Err(err))
			return nil, ewrap.ErrConnectionTime
		}
		return nil, ewrap.InternalError
	}

	return &server.RegisterResponse{
		UserId: userID,
	}, nil
}

func ValidateCredentials(req *server.RegisterRequest) error {
	// TODO: regex to check the email
	if req.GetEmail() == "" {
		return ewrap.ErrEmailRequired
	}
	// TODO: minimum password length
	if req.GetPassword() == "" {
		return ewrap.ErrPasswordRequired
	}
	return nil
}

func Login(ctx context.Context, req *server.LoginRequest) (*server.LoginResponse, error) {
	panic("not implemented")
}

func IsAdmin(ctx context.Context, req *server.AdminRequest) (*server.AdminResponse, error) {
	panic("not implemented")
}
