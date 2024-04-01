package service

import (
	"github.com/ikarizxc/http-server/internal/repository"
	"github.com/ikarizxc/http-server/internal/service/user"
)

type Service struct {
	user.IUserService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		user.NewUserService(repo),
	}
}
