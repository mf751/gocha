package main

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

type ClientList map[uuid.UUID]map[uuid.UUID]*Client

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	chatsID    []uuid.UUID
	userID     uuid.UUID

	egress chan Event
}

func newClient(
	conn *websocket.Conn,
	manager *Manager,
	userID uuid.UUID,
	chatsID []uuid.UUID,
) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
		userID:     userID,
		chatsID:    chatsID,
	}
}

func (c *Client) readMessages() {
	// cleanup function
	defer func() {
		c.manager.removeClient(c)
	}()

	// pinging
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		c.manager.app.logger.PrintError(
			err,
			map[string]string{"failed setting the deadline for the connection": err.Error()},
		)
		return
	}

	// jumbo frames
	c.connection.SetReadLimit(512)
	c.connection.SetPongHandler(c.pongHandler)

	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				c.manager.app.logger.PrintError(
					err,
					map[string]string{"error reading message: ": err.Error()},
				)
			}
			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			c.manager.app.logger.PrintError(
				err,
				map[string]string{"error unmarshaling event": err.Error()},
			)
			continue
		}

		if err := c.manager.routeEvents(request, c); err != nil {
			c.manager.app.logger.PrintError(
				err,
				map[string]string{"Error handling message": err.Error()},
			)
			continue
		}
	}
}

func (c *Client) pongHandler(pongMessage string) error {
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	ticker := time.NewTicker(pingInterval)

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					c.manager.app.logger.PrintError(
						err,
						map[string]string{"connection closed": err.Error()},
					)
				}
				break
			}
			data, err := json.Marshal(message)
			if err != nil {
				c.manager.app.logger.PrintError(
					err,
					map[string]string{"error marshaling data: ": err.Error()},
				)
				continue
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				c.manager.app.logger.PrintError(
					err,
					map[string]string{"failed to send message": err.Error()},
				)
				continue
			}
		// message sent
		case <-ticker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
				if errors.Is(err, websocket.ErrCloseSent) {
					break
				}
				c.manager.app.logger.PrintError(
					err,
					map[string]string{"writing ping message error": err.Error()},
				)
				return
			}
		}
	}
}
