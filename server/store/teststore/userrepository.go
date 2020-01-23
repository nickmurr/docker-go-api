package teststore

import (
	"github.com/nickmurr/go-http-rest-api/model"
	"github.com/nickmurr/go-http-rest-api/store"
)

// UserRepository ...
type UserRepository struct {
	store *Store
	users map[int]*model.User
}

// Create ...
func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	u.ID = len(r.users) + 1
	r.users[u.ID] = u

	return nil
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}

	return nil, store.RecordNotFound
}


func (r *UserRepository) FindById(id int) (*model.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, store.RecordNotFound
	}

	return u, nil
}
