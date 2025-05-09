package main

import (
	"sync"

	"github.com/google/uuid"
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
		clients:           make(ClientList),
		app:               app,
		connectionClients: make(map[uuid.UUID]*Client),
	}
	return m
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.connectionClients[client.userID] = client

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
		if len(m.clients[chatID]) == 0 {
			delete(m.clients, chatID)
		}
	}
}
