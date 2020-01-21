package store_test

import (
	"fmt"
	"github.com/nickmurr/go-http-rest-api/model"
	"github.com/nickmurr/go-http-rest-api/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	fmt.Println(databaseURL)
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("users")

	u, err := s.User().Create(&model.User{
		Email: "user@example.org",
	})
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
