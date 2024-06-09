package types

import (
	"time"
)

type User struct {
	ID 			int
	Username	string
	Email		string
	Password	string
	CreatedAt	time.Time
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) (int, error)
}