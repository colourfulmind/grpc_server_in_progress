package models

// User is a struct to store information about user
type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	PassHash []byte `json:"pass_hash"`
	IsAdmin  bool   `json:"is_admin"`
}

type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
