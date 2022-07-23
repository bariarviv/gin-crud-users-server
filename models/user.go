package models

import "time"

type User struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

// NewUser returns a new user
func NewUser(email, name, password string) *User {
	return &User{Email: email, Name: name, Password: password, CreatedAt: time.Now().String()}
}
