package main

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventNewMessage    string = "new_message"
	EventJoinedMessage string = "joined_message"
	EventLeftMessage   string = "left_message"
)

type NewMessageEvent struct {
	Message string    `json:"message"`
	ChatID  uuid.UUID `json:"chat_id"`
	Sent    time.Time `json:"sent"`
	From    uuid.UUID `json:"from"`
}
