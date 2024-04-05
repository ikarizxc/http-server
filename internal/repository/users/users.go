package users

import (
	"github.com/ikarizxc/http-server/internal/entities/users"
)

type Storage interface {
	Create(user *users.User) (int, error)
	GetById(id int) (*users.User, error)
	GetByEmail(email string) (*users.User, error)
	GetAll() ([]*users.User, error)
	Delete(id int) error
	Update(id int, fieldsToUpdate map[string]string) error
}
