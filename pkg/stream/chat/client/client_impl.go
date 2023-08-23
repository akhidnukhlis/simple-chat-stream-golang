package pkg

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	Time     string `json:"time"`
}

type ConsoleMessageHandler struct{}

func (c ConsoleMessageHandler) HandleMessage(msg Message) {
	fmt.Printf("[%s] %s: %s\n", msg.Time, msg.Username, msg.Text)
}

// SOLID Principle: Open/Closed Principle (OCP)
type WebSocketClient struct {
	conn      *websocket.Conn
	handler   MessageHandler
	interrupt chan os.Signal
	done      chan struct{}
}

func NewWebSocketClient(serverURL string, handler MessageHandler) (*WebSocketClient, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	return &WebSocketClient{
		conn:      conn,
		handler:   handler,
		interrupt: interrupt,
		done:      make(chan struct{}),
	}, nil
}

func (c *WebSocketClient) Start() {
	defer c.conn.Close()

	go c.readMessages()

	for {
		select {
		case <-c.done:
			return
		case <-c.interrupt:
			c.handleInterrupt()
			return
		}
	}
}

func (c *WebSocketClient) readMessages() {
	defer close(c.done)
	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Println("read:", err)
			return
		}
		c.handler.HandleMessage(msg)
	}
}

func (c *WebSocketClient) handleInterrupt() {
	log.Println("interrupt")
	err := c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
	select {
	case <-c.done:
	case <-time.After(time.Second):
	}
}
