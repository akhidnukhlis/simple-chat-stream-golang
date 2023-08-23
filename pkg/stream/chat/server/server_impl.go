package pkg

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	Time     string `json:"time"`
}

type TextMessageProcessor struct{}

func (t TextMessageProcessor) ProcessMessage(msg Message) {
	fmt.Printf("[%s] %s: %s\n", msg.Time, msg.Username, msg.Text)
}

// SOLID Principle: Open/Closed Principle (OCP)
type Chat struct {
	participants map[*websocket.Conn]bool
	broadcast    chan Message
	processor    MessageProcessor
}

func NewChat(processor MessageProcessor) *Chat {
	return &Chat{
		participants: make(map[*websocket.Conn]bool),
		broadcast:    make(chan Message),
		processor:    processor,
	}
}

func (c *Chat) SendMessage(msg Message) {
	c.broadcast <- msg
}

func (c *Chat) Start() {
	for {
		select {
		case msg := <-c.broadcast:
			c.processor.ProcessMessage(msg)
			for conn := range c.participants {
				conn.WriteJSON(msg)
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnection(chat *Chat, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	chat.participants[conn] = true

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			delete(chat.participants, conn)
			break
		}

		chat.SendMessage(Message{
			Username: msg.Username,
			Text:     msg.Text,
			Time:     time.Now().Format("15:04:05"),
		})
	}
}
