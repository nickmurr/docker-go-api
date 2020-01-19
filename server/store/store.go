package store

import (
	"database/sql"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Store struct {
	config *Config
	db     *sql.DB
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Open() error {
	url := fmt.Sprintf("postgresql://postgres:password@postgres:5432/example?sslmode=disable")

	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}

	// m, err := migrate.NewMigration()

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
