package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/mf751/gocha/internal/data"
)

var Upgrader = websocket.Upgrader{
	// CheckOrigin:     checkOrigin,
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (manager *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ChatID uuid.UUID `json:"chat_id"`
	}
	err := manager.app.readJSON(w, r, &input)
	if err != nil {
		manager.app.badRequestResponse(w, r, err)
	}
	user := manager.app.contextGetUser(r)
	err = manager.app.models.Users.IsInChat(user.ID, input.ChatID)
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
		manager.app.serverErrorResponse(w, r, err)
		return
	}

	client := newClient(conn, manager, user.ID, input.ChatID)

	manager.addClient(client)

	go client.readMessages(r)
	go client.writeMessages(r)
}

//
// func checkOrigin(r *http.Request) bool {
// 	origin := r.Header.Get("Origin")
//
// 	switch origin {
// 	case "http://localhost:5173":
// 		return true
// 	default:
// 		return false
// 	}
// }
