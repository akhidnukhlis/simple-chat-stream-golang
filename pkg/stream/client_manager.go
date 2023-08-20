package stream

type ClientManagerInterface interface {
	RegisterClient(client Client)
	UnregisterClient(client Client)
	BroadcastMessage(message string)
}
