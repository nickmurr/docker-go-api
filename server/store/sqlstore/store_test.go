package sqlstore_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	// 	...
	databaseURL = os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "user=postgres password=docker host=postgres dbname=restapi_test port=5432 sslmode=disable"
	}

	os.Exit(m.Run())
}
