package store_test

import (
	"fmt"
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	// 	...
	// databaseURL = os.Getenv("DATABASE_URL")
	databaseURL = "user=postgres password=docker host=postgres dbname=restapi_test port=5432 sslmode=disable"
	fmt.Println(databaseURL)
	if databaseURL == "" {
		databaseURL = "host=postgres dbname=restapi_test sslmode=disable"
	}

	os.Exit(m.Run())
}
