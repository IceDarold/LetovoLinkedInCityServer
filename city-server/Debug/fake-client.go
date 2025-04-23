package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"strconv"

	"github.com/gorilla/websocket"
)

type Vec3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type JoinMessage struct {
	PlayerID string `json:"playerId"`
}

type MoveMessage struct {
	PlayerID string `json:"playerId"`
	Position Vec3   `json:"position"`
}

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func sendJSON(conn *websocket.Conn, msgType string, data any) {
	payload, _ := json.Marshal(data)
	message := Message{
		Type: msgType,
		Data: payload,
	}
	full, _ := json.Marshal(message)
	conn.WriteMessage(websocket.TextMessage, full)
	fmt.Println("✅ Sent:", msgType)
}

func main() {
	// Connect
	playerID := "test-client-1"
	url := "ws://localhost:8080/ws"

	fmt.Printf("Connecting to %s as %s...\n", url, playerID)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("❌ Failed to connect:", err)
	}
	defer conn.Close()

	// Send join
	join := JoinMessage{PlayerID: playerID}
	sendJSON(conn, "player_joined", join)

	// Goroutine to listen for messages
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println("⚠️ Read error:", err)
				return
			}
			fmt.Println("📩 Server says:", string(msg))
		}
	}()

	// Input loop
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Command (move x y z / exit): ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if strings.HasPrefix(text, "move") {
			parts := strings.Split(text, " ")
			if len(parts) != 4 {
				fmt.Println("Usage: move x y z")
				continue
			}

			x := parseFloat(parts[1])
			y := parseFloat(parts[2])
			z := parseFloat(parts[3])

			move := MoveMessage{
				PlayerID: playerID,
				Position: Vec3{X: x, Y: y, Z: z},
			}

			sendJSON(conn, "player_moved", move)
		} else if text == "exit" {
			break
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func parseFloat(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Println("Invalid float:", s)
		return 0
	}
	return v
}
