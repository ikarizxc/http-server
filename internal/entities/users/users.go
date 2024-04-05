package users

type User struct {
	Id        int    `json:"-" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Username  string `json:"username" db:"username"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"-" db:"password_hash"`
	IsAdmin   bool   `json:"-" db:"is_admin"`
}
