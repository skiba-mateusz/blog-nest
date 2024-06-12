package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/skiba-mateusz/blog-nest/types"
)

type blogStore struct {
	db	*sql.DB
}

func NewBlogStore(db *sql.DB) *blogStore {
	return &blogStore{db: db}
}

func (s *blogStore) GetCategories() ([]types.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, `SELECT * FROM categories`)
	if err != nil {
		return nil, err
	}

	var categories []types.Category

	for rows.Next() {
		var category types.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (s *blogStore) GetBlogs() ([]types.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, `
		SELECT 
			b.id, b.title, b.content, b.created_at,
			u.id, u.username,
			c.id, c.name 
		FROM 
			blogs b
		INNER JOIN 
			users u on b.user_id = u.id
		INNER JOIN
			categories c on b.category_id = c.id	
	`)
	if err != nil {
		return nil, err
	}

	var blogs []types.Blog
	for rows.Next() {
		blog, err := scanRowsIntoBlog(rows)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, *blog)
	}

	return blogs, nil
}

func (s *blogStore) GetBlogByID(blogID int) (*types.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, `
		SELECT 
				b.id, b.title, b.content, b.created_at,
				u.id, u.username,
				c.id, c.name
		FROM blogs b
		JOIN users u ON b.user_id = u.id
		JOIN categories c ON b.category_id = c.id
		WHERE b.id = $1
	`, blogID)
	if err != nil {
		return nil, err
	}

	blog := new(types.Blog)
	for rows.Next() {
		blog, err = scanRowsIntoBlog(rows)
		if err != nil {
			return nil, err
		}
	}

	return blog, nil

}

func (s *blogStore) CreateBlog(blog types.Blog) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5) 
	defer cancel()

	var id int

	row := s.db.QueryRowContext(ctx, `INSERT INTO blogs (title, content, user_id, category_id) VALUES ($1, $2, $3, $4) returning id`,
		blog.Title,
		blog.Content,
		blog.User.ID,
		blog.Category.ID,
	)
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *blogStore) GetBlogLikes(userID, blogID int) (*types.BlogLikes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5) 
	defer cancel()

	blogLikes := new(types.BlogLikes)
	row := s.db.QueryRowContext(ctx, `
		SELECT 
			COALESCE(sum(value),0) AS likes_count, 
			CASE
				WHEN $2 = 0 THEN FALSE
				ELSE EXISTS(SELECT 1 from blog_likes WHERE blog_id = $1 AND user_id = $2)
			END AS user_liked,
			COALESCE((SELECT value FROM blog_likes WHERE blog_id = $1 AND user_id = $2),0) as user_like_value 
		FROM blog_likes 
		WHERE blog_id = $1
		`, blogID, userID,
	)
	err := row.Scan(&blogLikes.Count, &blogLikes.UserLiked, &blogLikes.UserLikeValue)
	if err != nil {
		return nil, err
	}

	return blogLikes, nil
}

func (s *blogStore) CreateLike(userID, blogID, value int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5) 
	defer cancel()

	_, err := s.db.ExecContext(ctx, `INSERT INTO blog_likes (value, blog_id, user_id) VALUES($1, $2, $3)`, value, blogID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *blogStore) UpdateLike(userID, blogID, value int) (error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5) 
	defer cancel()

	_, err := s.db.ExecContext(ctx, `UPDATE blog_likes SET value = $1 WHERE blog_id = $2 AND user_id = $3`, value, blogID, userID)
	if err != nil {
		return err
	}

	return nil
}

func scanRowsIntoBlog(rows *sql.Rows) (*types.Blog, error) {
	blog := new(types.Blog)
	err := rows.Scan(
		&blog.ID, 
		&blog.Title, 
		&blog.Content, 
		&blog.CreatedAt, 
		&blog.User.ID, 
		&blog.User.Username, 
		&blog.Category.ID,
		&blog.Category.Name,
	);
	if err != nil {
		return nil, err
	}

	return blog, nil
}