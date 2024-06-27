package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/skiba-mateusz/blog-nest/types"
)

type userStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *userStore {
	return &userStore{
		db: db,
	}
}

func (s *userStore) GetProfileByID(id int) (*types.Profile, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	query := `
		SELECT 
			u.username, u.bio, u.avatar_path, created_at,
			(SELECT COUNT(*) FROM comments WHERE user_id = $1) AS num_comments,
			(SELECT COUNT(*) FROM blogs WHERE user_id = $1) AS num_blog
		FROM 
			users u
		WHERE 
			id = $1
	`

	row := s.db.QueryRowContext(ctx, query, id)

	profile := new(types.Profile)
	err := row.Scan(
		&profile.Username, 
		&profile.Bio, 
		&profile.AvatarPath, 
		&profile.CreatedAt,
		&profile.NumComments,
		&profile.NumBogs,
	)
	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (s *userStore) GetUserByEmail(email string) (*types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	
	query := `SELECT id, username, email, password FROM users WHERE email = $1`
	row := s.db.QueryRowContext(ctx, query, email)

	user := &types.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

func (s *userStore) GetUserByID(id int) (*types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	
	query := `SELECT id, username, email, password, bio, avatar_path, created_at FROM users WHERE id = $1`
	row := s.db.QueryRowContext(ctx, query, id)

	user, err := scanRowIntoUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}


func (s *userStore) CreateUser(user types.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	stmt := `INSERT INTO users (username, email, password) VALUES($1, $2, $3) RETURNING id`

	row := s.db.QueryRowContext(ctx, stmt, user.Username, user.Email, user.Password)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %v", err)
	}

	return int(id), nil
}

func (s *userStore) UpdateUser(user types.User) (*types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	updates := []string{}
	args := []interface{}{}
	argID := 0

	if user.Username != "" {
		argID++
		updates = append(updates, fmt.Sprintf("username = $%d", argID))
		args = append(args, user.Username)
	}

	if user.Bio != "" {
		argID++
		updates = append(updates, fmt.Sprintf("bio = $%d", argID))
		args = append(args, user.Bio)
	}

	if user.AvatarPath != "" {
		argID++
		updates = append(updates, fmt.Sprintf("avatar_path = $%d", argID))
		args = append(args, user.AvatarPath)
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	argID++
	args = append(args, user.ID)
	stmt := fmt.Sprintf(`
		UPDATE users SET %s WHERE id = $%d RETURNING id, username, email, password, bio, avatar_path, created_at`,
		strings.Join(updates, ", "), argID,
	)
	row := s.db.QueryRowContext(ctx, stmt, args...)
	updatedUser, err := scanRowIntoUser(row)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func scanRowIntoUser(row *sql.Row) (*types.User, error) {
	user := new(types.User)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Bio, &user.AvatarPath, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}