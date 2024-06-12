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
	Likes		*BlogLikes
	CreatedAt	time.Time
}

type BlogLikes struct {
	Count			int
	UserLiked		bool 
	UserLikeValue 	int
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) (int, error)
}

type BlogStore interface {
	GetCategories() ([]Category, error)
	GetBlogs() ([]Blog, error)
	GetBlogByID(blogID int) (*Blog, error)
	CreateBlog(blog Blog) (int, error)
	CreateLike(userID, blogID, value int) (error)
	UpdateLike(userID, blogID, value int) (error)
	GetBlogLikes(userID, blogID int) (*BlogLikes, error) 
}