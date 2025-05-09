package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/mf751/gocha/internal/data"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin:     checkOrigin,
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

	chatID, err := uuid.Parse(r.URL.Query().Get("chat_id"))
	if err != nil {
		manager.app.errorResponse(w, r, http.StatusBadRequest, "chat_id not a valid uuid")
		return
	}

	err = manager.app.models.Users.IsInChat(user.ID, chatID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNotInChat):
			manager.app.errorResponse(w, r, http.StatusUnauthorized, err.Error())
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

	client := newClient(conn, manager, user.ID, chatID)

	manager.addClient(client)

	go client.readMessages()
	go client.writeMessages()
	fmt.Println(manager.clients)
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:5173":
		return true
	default:
		return false
	}
}
