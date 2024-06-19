package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/skiba-mateusz/blog-nest/types"
)

type commentStore struct {
	db *sql.DB
}

func NewCommentStore(db *sql.DB) *commentStore {
	return &commentStore{
		db: db,
	}
}

func (s commentStore) GetCommentsByBlogID(blogID int, userID int) ([]types.Comment, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5) 
	defer cancel()

	query := `
		SELECT 
			c.id, c.content, c.parent_id, c.created_at,
			b.id,
			u.id, u.username,
			COALESCE(SUM(cl.value), 0) AS likes_count,
			COALESCE((SELECT value FROM comment_likes WHERE user_id = $2 AND comment_id = c.id), 0) AS user_like_value,
			CASE
				WHEN $2 = 0 THEN FALSE
				ELSE EXISTS(SELECT 1 FROM comment_likes WHERE user_id = $2 AND comment_id =  c.id)
			END AS user_liked
		FROM 
			comments c
		INNER JOIN 
			users u ON u.id = c.user_id
		INNER JOIN 
			blogs b ON b.id = c.blog_id
		LEFT JOIN 
			comment_likes cl ON cl.comment_id = c.id
		WHERE 
			c.blog_id = $1
		GROUP BY 
			c.id, b.id, u.id
		ORDER BY 
			c.created_at DESC
	`

	rows, err := s.db.QueryContext(ctx, query, blogID, userID)	
	if err != nil {
		return nil, err
	}

	comments := []types.Comment{}

	for rows.Next() {
		comment, err := scanRowsIntoComment(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, *comment)
	}

	return comments, nil
}

func (s commentStore) CreateComment(comment types.Comment) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	stmt := `INSERT INTO comments (content, user_id, blog_id, parent_id) VALUES ($1, $2, $3, $4) RETURNING id`

	var commentID int
	row := s.db.QueryRowContext(ctx, stmt, comment.Content, comment.User.ID, comment.Blog.ID, comment.ParentID)

	err := row.Scan(&commentID)
	if err != nil {
		return 0, err
	}

	return commentID, nil
}

func (s commentStore) CreateLike(value, commentID, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	stmt := `INSERT INTO comment_likes (value, comment_id, user_id) VALUES($1, $2, $3)`

	_, err := s.db.ExecContext(ctx, stmt, value, commentID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s commentStore) UpdateLike(value, commentID, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	stmt := `UPDATE comment_likes SET value = $1 WHERE comment_id = $2 AND user_id = $3`
	
	_, err := s.db.ExecContext(ctx, stmt, value, commentID, userID)
	if err != nil {
		return err
	}
	
	return nil
}

func (s commentStore) GetCommentLikes(commentID, userID int) (*types.Likes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	query := `
		SELECT 
			COALESCE(SUM(cl.value), 0) AS likes_count,
			COALESCE((SELECT value FROM comment_likes WHERE user_id = $1 AND comment_id = $2), 0) AS user_like_value,
			CASE
	            WHEN 2 = 0 THEN FALSE
	            ELSE EXISTS(SELECT 1 FROM comment_likes WHERE comment_id = $2 AND user_id = $1)
	        END AS user_liked
		FROM 
			comment_likes cl
		WHERE 
			cl.comment_id = $2
	`

	likes := new(types.Likes)
	row := s.db.QueryRowContext(ctx, query, userID, commentID)
	
	err := row.Scan(&likes.Count, &likes.UserLikeValue, &likes.UserLiked)
	if err != nil {
		return nil, err
	}

	return likes, nil
}

func scanRowsIntoComment(rows *sql.Rows) (*types.Comment, error) {
	comment := new(types.Comment)

	err := rows.Scan(&comment.ID, 
		&comment.Content, 
		&comment.ParentID, 
		&comment.CreatedAt,
		&comment.Blog.ID, 
		&comment.User.ID, 
		&comment.User.Username,
		&comment.Likes.Count,
		&comment.Likes.UserLikeValue,
		&comment.Likes.UserLiked,
	)
	if err != nil {
		return nil, err
	}

	return comment, nil
}