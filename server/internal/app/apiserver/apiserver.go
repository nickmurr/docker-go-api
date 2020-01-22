package apiserver

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/nickmurr/go-http-rest-api/store/sqlstore"
	"net/http"
	"os"
)

func Start(config *Config) error {

	db, err := newDb(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte("secret"))
	s := newServer(store, sessionStore)
	fmt.Printf("Server running on port %v\n", os.Getenv("BIND_ADDR"))
	return http.ListenAndServe(":5000", s)
}

func newDb(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
