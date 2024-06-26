package postgres

import (
	"database/sql"
	"fmt"

	"github.com/ikarizxc/http-server/internal/entities/users"
	"github.com/ikarizxc/http-server/pkg/db/postgres"
	"github.com/jmoiron/sqlx"
)

type UsersStorage struct {
	db *sqlx.DB
}

func NewUsersStorage(db *sqlx.DB) *UsersStorage {
	return &UsersStorage{db: db}
}

func (r *UsersStorage) Create(user *users.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (first_name, last_name, username, email, password_hash) values ($1, $2, $3, $4, $5) RETURNING id", postgres.UsersTable)

	row := r.db.QueryRow(query, user.FirstName, user.LastName, user.Username, user.Email, user.Password)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UsersStorage) GetById(id int) (*users.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", postgres.UsersTable)

	var user users.User

	if err := r.db.Get(&user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user with id %d;", id)
		}
		return &user, err
	}

	return &user, nil
}

func (r *UsersStorage) GetByEmail(email string) (*users.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE email=$1;", postgres.UsersTable)

	var user users.User

	if err := r.db.Get(&user, query, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no user with email %s", email)
		}
		return &user, err
	}

	return &user, nil
}

func (r *UsersStorage) GetAll() ([]*users.User, error) {
	query := fmt.Sprintf("SELECT * FROM %s;", postgres.UsersTable)

	var users []*users.User

	if err := r.db.Select(&users, query); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UsersStorage) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", postgres.UsersTable)

	_, err := r.db.Exec(query, id)
	return err
}

func (r *UsersStorage) Update(id int, fieldsToUpdate map[string]string) error {
	query := fmt.Sprintf("UPDATE %s SET ", postgres.UsersTable)

	for k, v := range fieldsToUpdate {
		query += fmt.Sprintf("%s='%s' ", k, v)
	}
	query += fmt.Sprintf("WHERE id = %d;", id)

	_, err := r.db.Exec(query)
	return err
}
