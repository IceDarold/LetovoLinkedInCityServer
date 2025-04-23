package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Println("🔗 Attempting to upgrade to WebSocket...")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ Failed to upgrade:", err)
		return
	}

	log.Println("✅ WebSocket upgraded")

	client := &Client{
		Conn: conn,
		Send: make(chan []byte, 256),
		Hub:  hub,
	}
	log.Println("🌀 Client registration")

	hub.Register <- client
	log.Println("🌀 Starting ReadPump and WritePump")
	go client.WritePump()
	go client.ReadPump()
}
