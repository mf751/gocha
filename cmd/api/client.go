package main

import "github.com/gorilla/websocket"

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager

	egress chan Event
}
