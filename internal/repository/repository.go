package repository

import (
	"github.com/ikarizxc/http-server/internal/repository/user"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	user.IUserRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		user.NewUserRepositoryPostgres(db),
	}
}
