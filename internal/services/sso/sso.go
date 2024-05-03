package sso

import (
	"context"
	"errors"
	"log/slog"
	"main/internal/domain/models"
	"time"
)

var (
	ErrUserExists         = errors.New("user does not exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// SSO is a struct to authenticate users
type SSO struct {
	Log          *slog.Logger
	UserSaver    UserSaver
	UserProvider UserProvider
	TokenTTl     time.Duration
}

// UserSaver is an interface to save users to database
type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

// UserProvider is an interface that provides information about users
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}
