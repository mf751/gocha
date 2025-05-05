package main

import "sync"

type Manager struct {
	clients ClientList
	sync.RWMutex

	handlers map[string]EventHandler
}
