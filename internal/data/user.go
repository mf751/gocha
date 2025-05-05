package data

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	CreateAt  time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

var (
	ErrDuplicateEmial = errors.New("duplicate email")
	AnonymousUser     = &User{}
)

func (user *User) isAnonymous() bool {
	return user == AnonymousUser
}

type password struct {
	plainText *string
	hash      []byte
}

type UserModel struct {
	DB *sql.DB
}
