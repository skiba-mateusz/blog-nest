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