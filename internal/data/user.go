package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/mf751/gocha/internal/validator"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreateAt  Sent      `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated,omitempty"`
	// Version   int       `json:"-"`
}

var (
	ErrDuplicateEmial = errors.New("duplicate email")
	AnonymousUser     = &User{}
	ErrNotAdmin       = errors.New("User is not admin in chat")
)

func (user *User) IsAnonymous() bool {
	return user == AnonymousUser
}

type password struct {
	plainText *string
	hash      []byte
}

type UserModel struct {
	DB *sql.DB
}

func (psd *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	psd.plainText = &plainTextPassword
	psd.hash = hash

	return nil
}

func (psd *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(psd.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(vdtr *validator.Validator, email string) {
	vdtr.Check(
		validator.Matches(email, validator.EmailRX),
		"email",
		"must be a valid email address",
	)
}

func ValidatePasswordPlainText(vdtr *validator.Validator, password string) {
	vdtr.Check(password != "", "password", "must be provided")
	vdtr.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	vdtr.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(vdtr *validator.Validator, user *User) {
	vdtr.Check(user.Name != "", "name", "must be provided")
	vdtr.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(vdtr, user.Email)
	if user.Password.plainText != nil {
		ValidatePasswordPlainText(vdtr, *user.Password.plainText)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (model UserModel) Insert(user *User) error {
	sqlQuery := `
INSERT INTO users(id, name, email, password_hash, activated)
VALUES($1, $2, $3, $4, $5)
RETURNING id, created_at
	`
	args := []interface{}{
		uuid.New().String(),
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := model.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(&user.ID, &user.CreateAt.Sent)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmial
		default:
			return err
		}
	}
	return nil
}

func (model UserModel) GetByEmail(email string) (*User, error) {
	sqlQuery := `
SELECT id, created_at, name, email, password_hash, activated
FROM users
WHERE email = $1
	`
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := model.DB.QueryRowContext(ctx, sqlQuery, email).Scan(
		&user.ID,
		&user.CreateAt.Sent,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (model UserModel) Update(user *User) error {
	sqlQuery := `
UPDATE users
SET name = $1, email = $2, password_hash = $3, activated = $4
WHERE id = $5
	`
	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := model.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(&user.Name)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmial
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (model UserModel) GetForToken(tokenScope, tokenPlainText string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))

	sqlQuery := `
SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated
FROM users
INNER JOIN tokens
ON users.id = tokens.user_id
WHERE tokens.hash = $1
AND tokens.scope = $2
AND tokens.expiry > $3
	`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := model.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(
		&user.ID,
		&user.CreateAt.Sent,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (model UserModel) GetByID(ID uuid.UUID) (*User, error) {
	sqlQuery := `
SELECT name, created_at, password_hash, email, activated FROM users
WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	user := User{
		ID: ID,
	}
	err := model.DB.QueryRowContext(ctx, sqlQuery, ID).Scan(
		&user.Name,
		&user.CreateAt.Sent,
		&user.Password.hash,
		&user.Email,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (model UserModel) GetChats(UserID uuid.UUID) ([]*ChatWithLastMessage, error) {
	sqlQuery := `
WITH member_counts AS (
  SELECT chat_id, COUNT(*) AS member_count
  FROM users_chats
  GROUP BY chat_id
)
SELECT users.name,mc.member_count, chats.id, chats.name, chats.owner_id, chats.created_at, chats.is_private, messages.id, messages.sent, messages.user_id, messages.type, messages.content FROM users_chats
JOIN chats ON chats.id = users_chats.chat_id
JOIN users ON users.id = users_chats.user_id
JOIN member_counts mc ON mc.chat_id = chats.id
LEFT JOIN LATERAL(
	SELECT * FROM messages
	WHERE messages.chat_id = chats.id
	AND messages.deleted = false
	ORDER BY messages.sent DESC
	LIMIT 1
	) messages ON TRUE
WHERE users_chats.user_id = $1
ORDER BY messages.sent DESC NULLS LAST
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := model.DB.QueryContext(ctx, sqlQuery, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatsWithLastMessage []*ChatWithLastMessage

	for rows.Next() {
		var chatWithLastMessage ChatWithLastMessage
		err = rows.Scan(
			&chatWithLastMessage.LastMessage.User.Name,
			&chatWithLastMessage.Members,
			&chatWithLastMessage.Chat.ID,
			&chatWithLastMessage.Chat.Name,
			&chatWithLastMessage.Chat.OwnerID,
			&chatWithLastMessage.Chat.CreatedAt,
			&chatWithLastMessage.Chat.IsPrivate,
			&chatWithLastMessage.LastMessage.Message.ID,
			&chatWithLastMessage.LastMessage.Message.Sent.Sent,
			&chatWithLastMessage.LastMessage.Message.UserID,
			&chatWithLastMessage.LastMessage.Message.Type.Int,
			&chatWithLastMessage.LastMessage.Message.Content.NullString,
		)
		if err != nil {
			return nil, err
		}
		chatWithLastMessage.LastMessage.Message.ChatID = chatWithLastMessage.Chat.ID
		chatsWithLastMessage = append(chatsWithLastMessage, &chatWithLastMessage)
	}

	err = rows.Err()

	return chatsWithLastMessage, err
}

func (model UserModel) IsInChat(userID, chatID uuid.UUID) error {
	sqlQuery := `
SELECT TRUE FROM users_chats
WHERE user_id = $1
AND chat_id = $2
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var temp bool
	err := model.DB.QueryRowContext(ctx, sqlQuery, userID, chatID).Scan(&temp)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotInChat
	}

	return err
}

func (model UserModel) IsAdmin(userID, chatID uuid.UUID) error {
	sqlQuery := `
SELECT is_admin FROM users_chats
WHERE users_id = $1
AND chat_id = $2
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var isAdmin bool
	err := model.DB.QueryRowContext(ctx, sqlQuery, userID, chatID).Scan(&isAdmin)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotInChat
	}
	if !isAdmin {
		return ErrNotAdmin
	}

	return err
}

func (model UserModel) GetChatsID(UserID uuid.UUID) ([]uuid.UUID, error) {
	sqlQuery := `
SELECT chat_id FROM users_chats
WHERE user_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := model.DB.QueryContext(ctx, sqlQuery, UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatsID []uuid.UUID

	for rows.Next() {
		var chatID uuid.UUID
		err = rows.Scan(&chatID)
		if err != nil {
			return nil, err
		}
		chatsID = append(chatsID, chatID)
	}

	err = rows.Err()

	return chatsID, err
}
