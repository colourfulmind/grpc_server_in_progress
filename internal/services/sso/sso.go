package sso

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"main/internal/domain/models"
	"main/internal/storage"
	"main/pkg/jwt"
	"main/pkg/logger/sl"
	"time"
)

// TODO: add to yaml
const timeout = 3 * time.Second

var (
	ErrUserExists         = errors.New("user does not exists")
	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrConnectionTime     = errors.New("cannot connect to database")
)

// SSO is a struct to authenticate users
type SSO struct {
	Log      *slog.Logger
	Saver    UserSaver
	Provider UserProvider
	TokenTTl time.Duration
}

// UserSaver is an interface to save users to database
type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}

// UserProvider is an interface that provides information about users
type UserProvider interface {
	ProvideUser(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

// New return a new SSO structure
func New(log *slog.Logger, saver UserSaver, provider UserProvider, tokenTTL time.Duration) *SSO {
	return &SSO{
		Log:      log,
		Saver:    saver,
		Provider: provider,
		TokenTTl: tokenTTL,
	}
}

// RegisterNewUser generates a password hash and save the user to the database.
// Returns error if user already exists.
func (sso *SSO) RegisterNewUser(ctx context.Context, email, password string) (int64, error) {
	const op = "internal.services.sso.RegisterNewUser"

	log := sso.Log.With(slog.String("op", op), slog.String("email", email))
	log.Info("attempting to register user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	userID, err := sso.Saver.SaveUser(ctx, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", sl.Err(err))
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		if errors.Is(err, storage.ErrConnectionTime) {
			log.Error("connection time expired", sl.Err(err))
			return 0, fmt.Errorf("%s: %w", op, ErrConnectionTime)
		}
		log.Error("failed to register new user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user is registered", slog.Int64("user_id", userID))
	return userID, nil
}

// Login checks in the database if there is a user with the given email
// and then compare password to the password hash in the database.
// Returns error if user is not found else token.
func (sso *SSO) Login(ctx context.Context, email, password string, appID int32) (string, error) {
	const op = "internal.services.sso.Login"

	log := sso.Log.With(slog.String("op", op), slog.String("email", email))
	log.Info("attempting to login user")

	user, err := sso.Provider.ProvideUser(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		if errors.Is(err, storage.ErrConnectionTime) {
			log.Error("connection time expired", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrConnectionTime)
		}
		log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err = brypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token, err := jwt.New(user, sso.TokenTTl)
	if err != nil {
		log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user successfully logged in", slog.String("email", email))

	return token, nil
}

func (sso *SSO) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "internal.services.sso.IsAdmin"

	log := sso.Log.With(slog.String("op", op), slog.Int64("user_id", userID))
	log.Info("checking if user is admin")

	isAdmin, err := sso.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.Err(err))
			return false, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		if errors.Is(err, storage.ErrConnectionTime) {
			log.Error("connection time expired", sl.Err(err))
			return false, ErrConnectionTime
		}
		log.Error("failed to get user", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))
	return isAdmin, nil
}
