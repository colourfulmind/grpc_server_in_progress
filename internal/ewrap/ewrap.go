package ewrap

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	UserIdIsRequired      = status.Error(codes.InvalidArgument, "user id is required")
	ErrEmailRequired      = status.Error(codes.InvalidArgument, "email is required")
	ErrPasswordRequired   = status.Error(codes.InvalidArgument, "password is required")
	ErrInvalidCredentials = status.Error(codes.InvalidArgument, "incorrect email or password")
	UserAlreadyExists     = status.Error(codes.AlreadyExists, "user already exists")
	ErrConnectionTime     = status.Error(codes.DeadlineExceeded, "reached timeout while connecting to the database")
	InternalError         = status.Error(codes.Internal, "internal error")
)
