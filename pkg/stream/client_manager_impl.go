package stream

import (
	"sync"
)

type ChatClientManager struct {
	clients map[Client]bool
	mutex   sync.Mutex
}

func NewChatClientManager() *ChatClientManager {
	return &ChatClientManager{
		clients: make(map[Client]bool),
	}
}

func (cm *ChatClientManager) RegisterClient(client Client) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.clients[client] = true
	client.Start()
}

func (cm *ChatClientManager) UnregisterClient(client Client) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.clients, client)
}

func (cm *ChatClientManager) BroadcastMessage(message string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	for client := range cm.clients {
		client.SendMessage(message)
	}
}
