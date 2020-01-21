package store

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Store struct {
	config         *Config
	db             *sql.DB
	UserRepository *UserRepository
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Open() error {
	url := fmt.Sprintf(s.config.DatabaseURL)
	// url := fmt.Sprintf("postgresql://postgres:password@postgres:5432?sslmode=disable")

	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		fmt.Printf("Error while ping DB: %v", err)
		return err
	}

	s.db = db

	return nil
}

func (s *Store) Close() {
	_ = s.db.Close()
}

func (s *Store) User() *UserRepository {
	if s.UserRepository != nil {
		return s.UserRepository
	}

	s.UserRepository = &UserRepository{
		store: s,
	}

	return s.UserRepository
}
