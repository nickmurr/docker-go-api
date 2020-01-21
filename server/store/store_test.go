package store_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	// 	...
	databaseURL = os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// databaseURL = "postgres://postgres:docker@postgres:5432/restapi_test?sslmode=disable"
		// databaseURL = "user=postgres password=docker host=postgres dbname=restapi_test port=5432 sslmode=disable"
		databaseURL = "host=postgres dbname=restapi_test sslmode=disable"
	}

	os.Exit(m.Run())
}
