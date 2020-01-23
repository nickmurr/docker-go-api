package store

import "github.com/nickmurr/go-http-rest-api/model"

type UserRepository interface {
	Create(*model.User) error
	FindByEmail(string) (*model.User, error)
	FindById(int) (*model.User, error)
}
