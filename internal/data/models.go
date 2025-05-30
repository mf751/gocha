package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("records not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Modles struct {
	Users    UserModel
	Tokens   TokenModel
	Chats    ChatModel
	Messages MessagesModel
}

func NewModels(db *sql.DB) Modles {
	return Modles{
		Users:    UserModel{DB: db},
		Tokens:   TokenModel{DB: db},
		Chats:    ChatModel{DB: db},
		Messages: MessagesModel{DB: db},
	}
}
