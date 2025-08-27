package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Lobby struct {
	clients map[*websocket.Conn]bool
}

var lobbies = make(map[string]*Lobby)

func main() {
	// Создание нового лобби
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
			return
		}
		id := uuid.New().String()
		lobbies[id] = &Lobby{clients: make(map[*websocket.Conn]bool)}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"id":"%s"}`, id)
	})

	// WebSocket сигналинг
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		lobbyID := r.URL.Query().Get("lobbyID")
		lobby, ok := lobbies[lobbyID]
		if !ok {
			http.Error(w, "Lobby not found", http.StatusNotFound)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("WS upgrade error:", err)
			return
		}
		defer conn.Close()

		lobby.clients[conn] = true
		defer delete(lobby.clients, conn)

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("WS read error:", err)
				break
			}
			for client := range lobby.clients {
				if client != conn {
					client.WriteMessage(websocket.TextMessage, msg)
				}
			}
		}
	})

	// Главная страница — создание/ввод лобби
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "index.html")
			return
		}
		// Любой другой путь — это ID лобби
		http.ServeFile(w, r, "lobby.html")
	})

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
