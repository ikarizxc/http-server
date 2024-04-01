package user

type User struct {
	Id        int    `json:"-"`
	FirstName string `json:"first_name" db:"first_name" binding:"required"`
	LastName  string `json:"last_name" db:"last_name" binding:"required"`
	Username  string `json:"username" db:"username" binding:"required"`
	Email     string `json:"email" db:"email" binding:"required"`
	Password  string `json:"password" db:"password_hash" binding:"required"`
}
