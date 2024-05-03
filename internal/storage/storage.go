package storage

import "errors"

var (
	ErrUserExists     = errors.New("user already exists")
	ErrUserNotFound   = errors.New("user not found")
	ErrConnectionTime = errors.New("connection time is expired")

	//ErrAppNotFound  = errors.New("app not found")
)
