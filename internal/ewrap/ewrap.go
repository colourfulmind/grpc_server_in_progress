package ewrap

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrEmailRequired    = status.Error(codes.InvalidArgument, "email is required")
	ErrPasswordRequired = status.Error(codes.InvalidArgument, "password is required")
	UserAlreadyExists   = status.Error(codes.AlreadyExists, "user already exists")
	ErrConnectionTime   = status.Error(codes.DeadlineExceeded, "reached timeout while connecting to the database")
	InternalError       = status.Error(codes.Internal, "internal error")
)
