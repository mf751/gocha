package main

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/mf751/gocha/internal/data"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (manager *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	authToken := r.URL.Query().Get("token")
	user, err := manager.app.models.Users.GetForToken(data.ScopeAuthentication, authToken)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			manager.app.invalidAuthenticationTokenResponse(w, r)
			return
		default:
			manager.app.serverErrorResponse(w, r, err)
			return
		}
	}

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	chatsID, err := manager.app.models.Users.GetChatsID(user.ID)
	if err != nil {
		conn.Close()
		return
	}

	client := newClient(conn, manager, user.ID, chatsID)

	manager.addClient(client)

	go client.writeMessages()
}
