package model_test

import (
	"github.com/nickmurr/go-http-rest-api/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_BeforeCreate(t *testing.T) {
	u := model.TestUser(t)
	assert.NoError(t, u.BeforeCreate())
	assert.NotEmpty(t, u.EncryptedPassword)
}

type testCase struct {
	name    string
	u       func() *model.User
	isValid bool
}

func TestUser_Validate(t *testing.T) {
	// u := model.TestUser(t)
	// assert.NoError(t, u.Validate())

	testCases := []testCase{
		{
			name: "valid",
			u: func() *model.User {
				return model.TestUser(t)
			},
			isValid: true,
		},
		{
			name: "empty email",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = ""
				return u
			},
			isValid: false,
		},
		{
			name: "invalid email",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Email = "invalid"
				return u
			},
			isValid: false,
		},
		{
			name: "empty password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = ""
				return u
			},
			isValid: false,
		},
		{
			name: "short password",
			u: func() *model.User {
				u := model.TestUser(t)
				u.Password = "12"
				return u
			},
			isValid: false,
		},
		{
			name: "if encrypted password",
			u: func() *model.User {
				u:= model.TestUser(t)
				u.Password = ""
				u.EncryptedPassword = "encryptedpassword"
				return u
			},
			isValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.u().Validate())
			} else {
				assert.Error(t, tc.u().Validate())
			}
		})
	}
}
