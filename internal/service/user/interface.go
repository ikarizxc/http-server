package user

import "github.com/ikarizxc/http-server/internal/entities/user"

type IUserService interface {
	CreateUser(user *user.User) (int, error)
	GetUser(id int) (user.User, error)
	GetUsers() ([]*user.User, error)
	DeleteUser(id int) error
}
