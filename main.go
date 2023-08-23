package main

import (
	client "github.com/akhidnukhlis/simple-chat-stream-golang/pkg/stream/chat/client"
	server "github.com/akhidnukhlis/simple-chat-stream-golang/pkg/stream/chat/server"
	"log"
	"net/http"
	"net/url"
)

func main() {
	//start the server
	processor := server.TextMessageProcessor{}
	chat := server.NewChat(&processor)
	go chat.Start()

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		server.HandleConnection(chat, w, r)
	})

	if err := http.ListenAndServe(":8008", nil); err != nil {
		log.Printf("webSocket server got error: ", err)
		return
	}
	//http.ListenAndServe(":8008", nil)

	log.Printf("WebSocket server started on : 8008")

	//start the client
	u := url.URL{Scheme: "ws", Host: "localhost:8008", Path: "/chat"}
	handler := client.ConsoleMessageHandler{}
	socketClient, err := client.NewWebSocketClient(u.String(), handler)
	if err != nil {
		log.Printf("client creation error: ", err)
		return
	}

	socketClient.Start()
}
