package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/mf751/gocha/internal/data"
)

type Manager struct {
	clients           ClientList
	connectionClients map[uuid.UUID]*Client
	sync.RWMutex

	handlers map[string]EventHandler
	app      *application
}

func newManager(app *application) *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: map[string]EventHandler{},
		app:      app,
	}
	m.setupEventHandler()
	return m
}

func (m *Manager) setupEventHandler() {
	m.handlers[EventSendMessage] = sendMessage
}

func sendMessage(event Event, c *Client) error {
	var chatEvent SendMessageEvent

	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		return fmt.Errorf("bad payload in the request %v", err.Error())
	}

	err := c.manager.app.models.Users.IsInChat(c.userID, chatEvent.ChatID)
	if err != nil {
		return err
	}

	message := &data.Message{
		ID:     uuid.New(),
		UserID: c.userID,
		ChatID: chatEvent.ChatID,
		Content: sql.NullString{
			Valid:  true,
			String: chatEvent.Message,
		},
		Type: sql.NullInt32{
			Valid: true,
			Int32: data.MessageNormal,
		},
	}
	err = c.manager.app.models.Messages.SendMessage(message)
	if err != nil {
		return err
	}

	var broadMessage NewMessageEvent
	broadMessage.Message = message.Content.String
	broadMessage.From = message.UserID
	broadMessage.ChatID = message.ChatID
	broadMessage.Sent = message.Sent.Time

	sendData, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal the broadcast: %v", err)
	}

	outGoingEvent := Event{
		Payload: sendData,
		Type:    EventNewMessage,
	}

	for _, client := range c.manager.clients[chatEvent.ChatID] {
		client.egress <- outGoingEvent
	}

	return nil
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	for _, chatID := range client.chatsID {
		if _, ok := m.clients[chatID]; !ok {
			m.clients[chatID] = make(map[uuid.UUID]*Client)
		}
		m.clients[chatID][client.userID] = client
	}
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	client.connection.Close()
	delete(m.connectionClients, client.userID)
	for _, chatID := range client.chatsID {
		delete(m.clients[chatID], client.userID)
	}
}

func (m *Manager) routeEvents(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		return handler(event, c)
	}
	return errors.New("there is no such event type")
}
