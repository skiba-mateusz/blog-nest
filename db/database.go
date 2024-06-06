package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/skiba-mateusz/blog-nest/config"
)

const (
	maxOpenConns = 30
	maxIdleConns = 30
	maxIdleTime = 15 * time.Minute
)

func NewPostgreSQLStorage() (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", 
		config.Envs.DBUser,
		config.Envs.DBName,
		config.Envs.DBPassword,
	)
		
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleTime)

	if err = initStorage(db); err != nil {
		return nil, err
	}

	return db, nil
}

func initStorage(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	return nil
}