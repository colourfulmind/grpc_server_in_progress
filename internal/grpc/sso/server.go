package sso

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"main/internal/ewrap"
	"main/internal/services/sso"
	"main/pkg/logger/sl"
	server "main/protos/gen/go/blog"
	"regexp"
)

const EmptyValue = 0

// ServerSSO is a structure to store data for sso server describes in protobuf.
type ServerSSO struct {
	server.UnimplementedSSOServer
	sso SSO
	log *slog.Logger
}

// SSO is an interface which describes methods for the sso server.
type SSO interface {
	RegisterNewUser(ctx context.Context, email, password string) (int64, error)
	Login(ctx context.Context, email, password string, appID int32) (string, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

// Register registers a new sso server.
func Register(s *grpc.Server, sso SSO, log *slog.Logger) {
	server.RegisterSSOServer(s, &ServerSSO{
		sso: sso,
		log: log,
	})
}

// RegisterNewUser registers a new user. Returns errors in cases when the user already exists,
// connection to the database failed or incorrect email or password are provided.
func (s *ServerSSO) RegisterNewUser(ctx context.Context, req *server.RegisterRequest) (*server.RegisterResponse, error) {
	const op = "internal.grpc.sso.RegisterNewUser"

	log := s.log.With(slog.String("op", op), slog.String("email", req.GetEmail()))

	if err := ValidateCredentials(req.GetEmail(), req.GetPassword()); err != nil {
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
		log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, ewrap.InternalError
	}

	return &server.RegisterResponse{
		UserId: userID,
	}, nil
}

// Login logs in a user. Returns error if incorrect email or password are provided or
// connection to the database failed.
func (s *ServerSSO) Login(ctx context.Context, req *server.LoginRequest) (*server.LoginResponse, error) {
	const op = "internal.grpc.sso.Login"

	log := s.log.With(slog.String("op", op), slog.String("email", req.GetEmail()))

	if err := ValidateCredentials(req.GetEmail(), req.GetPassword()); err != nil {
		log.Error("invalid email or password", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	token, err := s.sso.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, sso.ErrInvalidCredentials) {
			log.Warn("invalid credentials", sl.Err(err))
			return nil, ewrap.ErrInvalidCredentials
		}
		if errors.Is(err, sso.ErrConnectionTime) {
			log.Error("failed to connect to the database", sl.Err(err))
			return nil, ewrap.ErrConnectionTime
		}
		log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, ewrap.InternalError
	}

	return &server.LoginResponse{
		Token: token,
	}, nil
}

// ValidateCredentials checks if email or password formats are correct.
func ValidateCredentials(email, password string) error {
	matched, err := regexp.MatchString(`([a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)`, email)
	if err != nil {
		return ewrap.ErrParsingRegex
	}
	if !matched || email == "" {
		return ewrap.ErrInvalidEmail
	}
	// TODO: minimum password length
	if password == "" {
		return ewrap.ErrPasswordRequired
	}
	return nil
}

// IsAdmin checks if user is an admin. Returns error if incorrect user id is provided.
func (s *ServerSSO) IsAdmin(ctx context.Context, req *server.AdminRequest) (*server.AdminResponse, error) {
	const op = "internal.grpc.sso.IsAdmin"

	log := s.log.With(slog.String("op", op), slog.Int64("user_id", req.GetUserId()))

	if err := ValidateAdmin(req.GetUserId()); err != nil {
		log.Warn("invalid credentials", sl.Err(err))
		return nil, ewrap.ErrInvalidCredentials
	}

	isAdmin, err := s.sso.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, sso.ErrInvalidCredentials) {
			log.Warn("invalid credentials", sl.Err(err))
			return nil, ewrap.ErrInvalidCredentials
		}
		if errors.Is(err, sso.ErrConnectionTime) {
			log.Error("failed to connect to the database", sl.Err(err))
			return nil, ewrap.ErrConnectionTime
		}
		log.Error("internal error", slog.String("op", op), sl.Err(err))
		return nil, ewrap.InternalError
	}

	return &server.AdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

// ValidateAdmin checks if provided user id is correct.
func ValidateAdmin(userID int64) error {
	if userID == EmptyValue {
		return ewrap.UserIdIsRequired
	}
	return nil
}
