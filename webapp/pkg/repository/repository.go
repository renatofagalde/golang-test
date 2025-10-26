package repository

import (
	"bootstrap/pkg/data"
	"database/sql"
)

type DatabaseRepository interface {
	Connection() *sql.DB
	GetUser(id int) (*data.User, error)
	GetUserByEmail(email string) (*data.User, error)
	AllUsers() ([]*data.User, error)
	UpdateUser(u data.User) error
	DeleteUser(id int) error
	InsertUser(user data.User) (int, error)
	ResetPassword(id int, password string) error
	InsertUserImage(i data.UserImage) (int, error)
}
