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

type Category struct {
	ID		int
	Name	string
}

type Blog struct {
	ID			int
	Title		string
	Content		string
	Category	Category
	User		User
	CreatedAt	time.Time
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) (int, error)
}

type BlogStore interface {
	GetCategories() ([]Category, error)
	CreateBlog(blog Blog) (int, error)
}