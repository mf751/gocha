package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"

	"github.com/mf751/gocha/internal/data"
)

type Manager struct {
	clients ClientList
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
		return fmt.Errorf("bad payload in the request %v", err)
	}

	message := &data.Message{
		ID:     uuid.New(),
		UserID: c.userID,
		ChatID: c.chatID,
	}
	message.Content.String = chatEvent.Message
	message.Type.Int32 = 1
	err := c.manager.app.models.Messages.SendMessage(message)
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

	for client := range c.manager.clients[c.chatID] {
		client.egress <- outGoingEvent
	}

	return nil
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client.chatID] = make(map[*Client]bool)
	m.clients[client.chatID][client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client.chatID][client]; ok {
		client.connection.Close()
		delete(m.clients[client.chatID], client)
	}
}

func (m *Manager) routeEvents(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		return handler(event, c)
	}
	return errors.New("there is no such event type")
}
