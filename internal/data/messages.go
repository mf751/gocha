package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/mf751/gocha/internal/validator"
)

type Message struct {
	ID      uuid.UUID `json:"id"`
	Sent    Sent      `json:"sent"`
	ChatID  uuid.UUID `json:"chat_id"`
	UserID  uuid.UUID `json:"user_id"`
	Content Content   `json:"content"`
	Type    Int32     `json:"type"`
}

type Content struct {
	NullString sql.NullString
}

func (cnt Content) MarshalJSON() ([]byte, error) {
	return json.Marshal(cnt.NullString.String)
}

type Sent struct {
	Sent sql.NullTime
}

func (s Sent) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Sent.Time)
}

type Int32 struct {
	Int sql.NullInt32
}

func (i Int32) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.Int.Int32)
}

type MessagesModel struct {
	DB *sql.DB
}

const (
	MessageJoined = int32(50)
	MessageLeft   = int32(51)
	MessageNormal = int32(1)
)

var ErrMessageDeletionFailed = errors.New("failed to delete message")

func ValidateMessage(vdtr *validator.Validator, message *Message, userModel *UserModel) {
	vdtr.Check(message.Content.NullString.String != "", "content", "cannot be empty")
	vdtr.Check(
		len(message.Content.NullString.String) < 500,
		"content",
		"cannot be more than 500 characters long")
	vdtr.Check(
		message.Type.Int.Int32 == MessageNormal,
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

	args := []interface{}{
		message.ID,
		message.ChatID,
		message.UserID,
		message.Content.NullString,
		message.Type.Int,
	}

	err := model.DB.QueryRowContext(ctx, sqlQuery, args...).Scan(&message.Sent.Sent)

	return err
}

func (model MessagesModel) DeleteMessage(messageId, userID uuid.UUID, isAdmin bool) error {
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
