package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

type InputMessage struct {
	PlayerID   string  `json:"playerId"`
	Horizontal float64 `json:"horizontal"`
	Vertical   float64 `json:"vertical"`
	Jump       bool    `json:"jump"`
	Sprint     bool    `json:"sprint"`
	Dance      bool    `json:"dance"`
}

type JoinMessage struct {
	PlayerID string `json:"playerId"`
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
	playerID := "test-client-1"
	url := "ws://veconomics.ru:8082/ws"

	fmt.Printf("Connecting to %s as %s...\n", url, playerID)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("❌ Failed to connect:", err)
	}
	defer conn.Close()

	// Send join
	join := JoinMessage{PlayerID: playerID}
	sendJSON(conn, "player_joined", join)

	// Read messages
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

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Command (input h v [jump] [sprint] [dance] / exit): ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if strings.HasPrefix(text, "input") {
			parts := strings.Split(text, " ")
			if len(parts) < 3 {
				fmt.Println("Usage: input horizontal vertical [jump] [sprint] [dance]")
				continue
			}

			h := parseFloat(parts[1])
			v := parseFloat(parts[2])
			jump := parseBoolSafe(parts, 3)
			sprint := parseBoolSafe(parts, 4)
			dance := parseBoolSafe(parts, 5)

			input := InputMessage{
				PlayerID:   playerID,
				Horizontal: h,
				Vertical:   v,
				Jump:       jump,
				Sprint:     sprint,
				Dance:      dance,
			}

			sendJSON(conn, "player_input", input)

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

func parseBoolSafe(parts []string, index int) bool {
	if len(parts) > index {
		return parts[index] == "true" || parts[index] == "1"
	}
	return false
}
