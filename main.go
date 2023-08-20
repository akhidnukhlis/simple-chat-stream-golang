package main

import (
	"fmt"
	"github.com/akhidnukhlis/simple-chat-stream-golang/pkg/stream"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleConnections(cMgr stream.ClientManagerInterface, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client := stream.NewChatClient(conn, cMgr)
	cMgr.RegisterClient(client)
}

func main() {
	clientManager := stream.NewChatClientManager()

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(clientManager, w, r)
	})

	fmt.Println("Server chat berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
