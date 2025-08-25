package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Failed to read:", err)
			break
		}

		fmt.Printf("New message: %s\n", msg)

		err = conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Failed to send:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
