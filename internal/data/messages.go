package data

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"

	"github.com/mf751/gocha/internal/validator"
)

type Message struct {
	ID      uuid.UUID      `json:"id"`
	Sent    sql.NullTime   `json:"sent"`
	ChatID  uuid.UUID      `json:"chat_id"`
	UserID  uuid.UUID      `json:"user_id"`
	Content sql.NullString `json:"content"`
	Type    sql.NullInt32  `json:"type"`
}

type MessagesModel struct {
	DB *sql.DB
}

var ErrMessageDeletionFailed = errors.New("failed to delete message")

func ValidateMessage(vdtr *validator.Validator, message *Message, userModel *UserModel) {
	vdtr.Check(message.Content.String == "", "content", "cannot be empty")
	vdtr.Check(
		len(message.Content.String) < 500,
		"content",
		"cannot be more than 500 characters long")
	vdtr.Check(
		errors.Is(userModel.IsInChat(message.UserID, message.ChatID), ErrNotInChat),
		"user",
		"must be a member of chat",
	)
	vdtr.Check(
		slices.Contains([]int32{1, 50, 51}, message.Type.Int32),
		"type",
		"unsupported message type",
	)
}

func (model MessagesModel) SendMessage(message *Message) error {
	sqlQuery := `
INSERT INTO messages(id, chat_id, user_id, content, type)
VALUES ($1, $2, $3, $4, $5)
RETURNING sent
	`
	message.ID = uuid.New()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{message.ID, message.ChatID, message.UserID, message.Content, message.Type}

	err := model.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(&message.Sent)

	return err
}

func (model MessagesModel) deleteMessage(messageId, userID uuid.UUID, isAdmin bool) error {
	sqlQuery := `
UPDATE messages
SET deleted = true
WHERE id = $1
AND ( messages.user_id = $2 OR $3 = true)
RETURNING true
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var temp bool
	err := model.DB.QueryRowContext(ctx, sqlQuery, messageId, userID, isAdmin).Scan(&temp)
	if err != sql.ErrNoRows {
		return ErrMessageDeletionFailed
	}
	return err
}
