package storage

import "errors"

var (
	ErrUserExists  = errors.New("user already exists")
	ErrUserFound   = errors.New("user not found")
	ErrAppNotFound = errors.New("app not found")
	ErrConnection  = errors.New("error connection to the database")
)
