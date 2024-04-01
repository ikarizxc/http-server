package user

import (
	"fmt"

	"github.com/ikarizxc/http-server/internal/db/postgres"
	"github.com/ikarizxc/http-server/internal/entities/user"
	"github.com/jmoiron/sqlx"
)

type UserRepositoryPostgres struct {
	db *sqlx.DB
}

func NewUserRepositoryPostgres(db *sqlx.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{db: db}
}

func (r *UserRepositoryPostgres) CreateUser(user *user.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (first_name, last_name, username, email, password_hash) values ($1, $2, $3, $4, $5) RETURNING id", postgres.UsersTable)

	row := r.db.QueryRow(query, user.FirstName, user.LastName, user.Username, user.Email, user.Password)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepositoryPostgres) GetUser(id int) (*user.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", postgres.UsersTable)

	var user *user.User

	if err := r.db.Get(&user, query, id); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepositoryPostgres) GetUsers() ([]*user.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s", postgres.UsersTable)

	var users []*user.User

	if err := r.db.Select(&users, query); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoryPostgres) DeleteUser(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", postgres.UsersTable)

	_, err := r.db.Exec(query, id)
	return err
}
