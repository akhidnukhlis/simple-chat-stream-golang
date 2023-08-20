package stream

type Client interface {
	Start()
	SendMessage(message string)
	Close()
}
