package user

import (
	"crypto/sha1"
	"fmt"

	"github.com/ikarizxc/http-server/internal/entities/user"
	"github.com/ikarizxc/http-server/internal/repository"
)

const (
	salt = "fmefjjfmeimf2p3432fr3"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *user.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *UserService) GetUser(id int) (user.User, error) {
	return s.repo.GetUser(id)
}

func (s *UserService) GetUsers() ([]*user.User, error) {
	return s.repo.GetUsers()
}

func (s *UserService) DeleteUser(id int) error {
	return s.repo.DeleteUser(id)
}

func (s *UserService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
