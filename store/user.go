package store

import (
	"context"
	"database/sql"
	"fmt"
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

func (s *userStore) GetUserByEmail(email string) (*types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = $1`
	row := s.db.QueryRowContext(ctx, query, email)

	user, err := scanRowIntoUser(row)
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
	
	query := `SELECT id, username, email, password, created_at FROM users WHERE id = $1`
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

func scanRowIntoUser(row *sql.Row) (*types.User, error) {
	user := new(types.User)
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}