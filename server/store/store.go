package store

import (
	"database/sql"
	"fmt"
	"github.com/nickmurr/go-http-rest-api/model"

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

	c := make(chan *model.User)
	go func() {

		user, err := s.User().Create(&model.User{
			Email:             "mail@gmail.com",
			EncryptedPassword: "1234",
		})

		if err != nil {
			fmt.Println(err)
		}

		c <- user

	}()
	fmt.Println(<-c)

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
