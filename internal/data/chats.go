package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/mf751/gocha/internal/validator"
)

var (
	ErrDuplicateChat  = errors.New("Duplicate chat with the same person")
	ErrDeletionFailed = errors.New("Failed to delete chat")
	ErrChatNotFound   = errors.New("Chat not found")
	ErrPrivateChat    = errors.New("Private Chat")
	ErrAlreadyMember  = errors.New("Already a member of chat")
	ErrNotInChat      = errors.New("Not a membor of chat")
)

type ChatModel struct {
	DB *sql.DB
}

type Chat struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	OwnerID   uuid.UUID `json:"owner_id"`
	IsPrivate bool      `json:"is_private"`
}

type ChatUser struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	IsAdmin bool      `json:"is_admin"`
}

type ChatWithLastMessage struct {
	Chat        Chat    `json:"chat"`
	LastMessage Message `json:"last_message"`
	Members     int     `json:"members"`
}

func ValidateChatName(vdtr *validator.Validator, name string) {
	vdtr.Check(len(name) <= 50, "name", "cannot be more than 50 characters long")
	vdtr.Check(name != "", "name", "must be provided")
}

func (model ChatModel) Insert(chat *Chat) (time.Time, error) {
	sqlQuery := `
INSERT INTO chats(id, name, owner_id, is_private)
VALUES( $1, $2, $3, $4)
RETURNING created_at
	`
	sqlQuery2 := `
INSERT INTO users_chats(user_id, chat_id, is_admin)
VALUES($1, $2, true)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{chat.ID, chat.Name, chat.OwnerID, chat.IsPrivate}
	err := model.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(&chat.CreatedAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "chats_pkey"`:
			return time.Time{}, ErrDuplicateChat
		default:
			return time.Time{}, err
		}
	}
	if chat.IsPrivate {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		model.DB.QueryRowContext(ctx, sqlQuery2, chat.OwnerID, chat.ID)
	}
	return chat.CreatedAt, nil
}

func (model ChatModel) Delete(userID, chatID uuid.UUID) error {
	sqlQuery := `
DELETE FROM chats 
WHERE id = $1
AND owner_id = $2
RETURNING id
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := model.DB.QueryRowContext(ctx, sqlQuery, chatID, userID).Scan(&chatID)
	if err == sql.ErrNoRows {
		return ErrDeletionFailed
	}
	return err
}

func (model ChatModel) GetChat(chat *Chat) error {
	sqlQuery := `
SELECT name, owner_id, created_at, is_private FROM chats 
WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := model.DB.QueryRowContext(ctx, sqlQuery, chat.ID).Scan(
		&chat.Name,
		&chat.OwnerID,
		&chat.CreatedAt,
		&chat.IsPrivate,
	)

	if err == sql.ErrNoRows {
		return ErrChatNotFound
	}
	return err
}

func (model ChatModel) GetUsers(chatID uuid.UUID) ([]*ChatUser, error) {
	chat := Chat{ID: chatID}
	err := model.GetChat(&chat)
	if err != nil {
		return nil, err
	}

	if chat.IsPrivate {
		return nil, ErrPrivateChat
	}

	sqlQuery := `
SELECT users.id, users.name, users_chats.is_admin FROM users_chats
JOIN users ON users.id = users_chats.user_id
WHERE users_chats.chat_id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := model.DB.QueryContext(ctx, sqlQuery, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatUsers []*ChatUser

	for rows.Next() {
		var user ChatUser
		err = rows.Scan(&user.ID, &user.Name, &user.IsAdmin)
		if err != nil {
			return nil, err
		}
		chatUsers = append(chatUsers, &user)
	}

	err = rows.Err()

	return chatUsers, err
}

func (model ChatModel) Join(chatID, userID uuid.UUID, isAdmin bool) error {
	chat := Chat{ID: chatID}
	err := model.GetChat(&chat)
	if err != nil {
		return err
	}

	if chat.IsPrivate {
		return ErrPrivateChat
	}

	sqlQuery := `
INSERT INTO users_chats(user_id, chat_id, is_admin)
VALUES($1, $2, $3)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = model.DB.ExecContext(ctx, sqlQuery, userID, chatID, isAdmin)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_chats_pkey"`:
			return ErrAlreadyMember
		default:
			return err
		}
	}
	return nil
}

func (model ChatModel) Leave(chatID, userID uuid.UUID) error {
	sqlQuery := `
DELETE FROM users_chats
WHERE user_id = $1
AND chat_id = $2
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := model.DB.ExecContext(ctx, sqlQuery, userID, chatID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotInChat
	}
	return nil
}

func (model ChatModel) GetChatMessage(chatID uuid.UUID, size, start int) ([]*Message, error) {
	sqlQuery := `
SELECT id, user_id, content, sent, type FROM messages
WHERE chat_id = $1
AND deleted = false
ORDER BY sent DESC
LIMIT $2
OFFSET $3
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := model.DB.QueryContext(ctx, sqlQuery, chatID, size, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message

	for rows.Next() {
		var message Message
		err = rows.Scan(
			&message.ID,
			&message.UserID,
			&message.Content.NullString,
			&message.Sent.Sent,
			&message.Type.Int,
		)
		if err != nil {
			return nil, err
		}
		message.ChatID = chatID
		messages = append(messages, &message)
	}

	err = rows.Err()
	return messages, err
}
