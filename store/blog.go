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

	query := `SELECT * FROM categories`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	 categories := []types.Category{}

	for rows.Next() {
		category := types.Category{}
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (s *blogStore) GetBlogs(offset, limit int, searchQuery, category string) ([]types.Blog, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	var totalBlogs int	
	blogs := []types.Blog{}

	row := s.db.QueryRowContext(ctx, `
		SELECT 
			count(*) 
		FROM 
			blogs b 
		INNER JOIN
			categories c ON b.category_id = c.id
		WHERE 
			($1 = '' OR b.title ILIKE '%' || $1 || '%')
			AND ($2 = '' OR c.name = $2)
			AND ($1 != '' OR b.id NOT IN (
				SELECT id
				FROM blogs 
				ORDER BY created_at DESC 
				LIMIT 4
			))
		`,	 searchQuery, category)
	if err := row.Scan(&totalBlogs); err != nil {
		return blogs, 0, err
	}

	query := `
		SELECT 
			b.id, b.title, c.name 
		FROM 
			blogs b
		INNER JOIN
			categories c on b.category_id = c.id
		WHERE
			($2 = '' OR b.title ILIKE '%' || $2 || '%') 
			AND ($3 = '' OR c.name = $3)
			AND ($2 != '' OR b.id NOT IN (
				SELECT id
				FROM blogs 
				ORDER BY created_at DESC 
				LIMIT 4
			))
		ORDER BY 
			b.created_at DESC
		LIMIT 
			$4
		OFFSET
			$1
	`

	rows, err := s.db.QueryContext(ctx, query, offset, searchQuery, category, limit)
	if err != nil {
		return nil, 0, err
	}

	for rows.Next() {
		blog := types.Blog{}
		err := rows.Scan(&blog.ID, &blog.Title, &blog.Category.Name)
		if err != nil {
			return nil, 0, err
		}
		blogs = append(blogs, blog)
	}

	return blogs, totalBlogs, nil
}

func (s *blogStore) GetLatestBlogs() ([]types.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	query := `
		SELECT 
			b.id, b.title,
			c.name 
		FROM 
			blogs b
		INNER JOIN
			categories c ON c.id = b.category_id
		ORDER BY 
			b.created_at DESC 
		LIMIT 4
	`
	blogs := []types.Blog{}

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return blogs, err
	}

	for rows.Next() {
		blog := types.Blog{}
		err := rows.Scan(&blog.ID, &blog.Title, &blog.Category.Name)
		if err != nil {
			return blogs, err
		}
		blogs = append(blogs, blog)
	}

	return blogs, nil
}

func (s *blogStore) GetBlogByID(blogID, userID int) (*types.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	query := `
		SELECT 
			b.id, b.title, b.content, b.created_at,
			u.id, u.username,
			c.id, c.name,
			COALESCE(SUM(bl.value),0) AS likes_count,
			COALESCE((SELECT value FROM blog_likes WHERE blog_id = $1 AND user_id = $2), 0) AS user_like_value,
			CASE
			    WHEN 1 = 0 THEN FALSE
			    ELSE EXISTS(SELECT 1 FROM blog_likes WHERE user_id = $2 AND blog_id = b.id)
			END AS user_liked
		FROM 
			blogs b
		INNER JOIN 
			users u ON b.user_id = u.id
		INNER JOIN 
			categories c ON b.category_id = c.id
		LEFT JOIN
			blog_likes bl ON b.id = bl.blog_id
		WHERE 
			b.id = $1
		GROUP BY 
		    u.id, b.id, c.id
	`

	rows, err := s.db.QueryContext(ctx, query, blogID, userID)
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

	stmt := `
		INSERT INTO 
			blogs (title, content, user_id, category_id) 
		VALUES 
			($1, $2, $3, $4) 
		RETURNING 
			id
	`

	row := s.db.QueryRowContext(
		ctx, 
		stmt,
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

func (s *blogStore) GetBlogLikes(userID, blogID int) (*types.Likes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5) 
	defer cancel()

	query := `
		SELECT 
			COALESCE(sum(value),0) AS likes_count, 
			CASE
				WHEN $2 = 0 THEN FALSE
				ELSE EXISTS(SELECT 1 from blog_likes WHERE blog_id = $1 AND user_id = $2)
			END AS user_liked,
			COALESCE((SELECT value FROM blog_likes WHERE blog_id = $1 AND user_id = $2),0) as user_like_value 
		FROM
			blog_likes 
		WHERE 
			blog_id = $1
	`

	likes := new(types.Likes)
	row := s.db.QueryRowContext(ctx, query, blogID, userID)
	err := row.Scan(&likes.Count, &likes.UserLiked, &likes.UserLikeValue)
	if err != nil {
		return nil, err
	}

	return likes, nil
}

func (s *blogStore) CreateLike(userID, blogID, value int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5) 
	defer cancel()

	stmt := `
		INSERT INTO 
			blog_likes (value, blog_id, user_id) 
		VALUES
			($1, $2, $3)
	`

	_, err := s.db.ExecContext(ctx, stmt, value, blogID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *blogStore) UpdateLike(userID, blogID, value int) (error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5) 
	defer cancel()
	
	stmt := `
		UPDATE 
			blog_likes 
		SET 
			value = $1 
		WHERE 
			blog_id = $2 AND user_id = $3
	`

	_, err := s.db.ExecContext(ctx, stmt, value, blogID, userID)
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
		&blog.Likes.Count,
		&blog.Likes.UserLikeValue,
		&blog.Likes.UserLiked,
	);
	if err != nil {
		return nil, err
	}

	return blog, nil
}