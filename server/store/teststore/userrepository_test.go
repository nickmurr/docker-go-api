package teststore_test

import (
	"github.com/nickmurr/go-http-rest-api/model"
	"github.com/nickmurr/go-http-rest-api/store"
	"github.com/nickmurr/go-http-rest-api/store/teststore"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststore.New()
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.NotNil(t, u)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststore.New()
	email := "user@example.org"
	_, err := s.User().FindByEmail(email)
	assert.EqualError(t, err, store.RecordNotFound.Error())

	u := model.TestUser(t)
	u.Email = email

	_ = s.User().Create(u)

	u, err = s.User().FindByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}
