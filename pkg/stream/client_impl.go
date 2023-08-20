package stream

import (
	"github.com/gorilla/websocket"
	"log"
)

type ChatClient struct {
	conn *websocket.Conn
	cMgr ClientManagerInterface
}

func NewChatClient(conn *websocket.Conn, cMgr ClientManagerInterface) *ChatClient {
	return &ChatClient{
		conn: conn,
		cMgr: cMgr,
	}
}

func (c *ChatClient) Start() {
	go c.receiveMessages()
}

func (c *ChatClient) receiveMessages() {
	defer c.Close()
	for {
		_, msgBytes, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		msg := string(msgBytes)
		c.cMgr.BroadcastMessage(msg)
	}
}

func (c *ChatClient) SendMessage(message string) {
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Println("write:", err)
	}
}

func (c *ChatClient) Close() {
	c.conn.Close()
	c.cMgr.UnregisterClient(c)
}
