package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// Check origin for security in production
	CheckOrigin: func(r *http.Request) bool { return true },
}

// TODO: implement websocket for realtime update

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer conn.Close()
	if err := conn.WriteMessage(websocket.TextMessage, []byte("connected successfully")); err != nil {
		fmt.Println("Error sending message")
	}

	fmt.Println("user connected")
	// for message := range dataChannel {
	//     err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	//     if err != nil {
	//         fmt.Println("Write error:", err)
	//         break
	//     }
	// }
}
