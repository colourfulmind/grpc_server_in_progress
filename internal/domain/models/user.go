package models

// User is a struct to store information about user
type User struct {
	ID       int64
	Email    string
	PassHash []byte
	IsAdmin  bool
}
