package types

import (
	"io"
	"time"
)

type User struct {
	ID 			int
	Username	string
	Email		string
	Password	string
	Bio			string
	AvatarPath	string
	CreatedAt	time.Time
}

type Profile struct {
	Username	string
	Bio			string
	AvatarPath	string
	NumComments int
	NumBogs		int
	IsOwner 	bool
	CreatedAt	time.Time
}

type Category struct {
	ID		int
	Name	string
}

type Blog struct {
	ID				int
	Title			string
	Content			string
	ThumbnailPath 	string
	Category		Category
	User			User
	Likes			Likes
	CreatedAt		time.Time
}

type Comment struct {
	ID			int
	Content		string
	ParentID	int
	User 		User
	Blog		Blog
	Likes 		Likes
	Replies 	[]Comment
	CreatedAt	time.Time
}

type Likes struct {
	Count			int
	UserLiked		bool 
	UserLikeValue 	int
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	GetProfileByID(id int) (*Profile, error)
	CreateUser(user User) (int, error)
	UpdateUser(user User) (*User, error)
}

type BlogStore interface {
	GetCategories() ([]Category, error)
	GetBlogs(offset, limit int, searchQuery, category string) ([]Blog, int, error)
	GetLatestBlogs() ([]Blog, error)
	GetBlogByID(blogID, userID int) (*Blog, error)
	CreateBlog(blog Blog) (int, error)
	CreateLike(userID, blogID, value int) (error)
	UpdateLike(userID, blogID, value int) (error)
	GetBlogLikes(userID, blogID int) (*Likes, error) 
}

type CommentStore interface {
	CreateComment(comment Comment) (int, error)
	GetCommentsByBlogID(blogID, userID int) ([]Comment, error)
	GetCommentLikes(commentID, userID int) (*Likes, error)
	CreateLike(value, commentID, userID int) (error)
	UpdateLike(value, commentID, userID int) (error)
}

type S3Uploader interface {
	PutObject(file io.Reader, filename, directory string) (string, error)
}