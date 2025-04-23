package ws

import (
	"encoding/json"
	"log"
	"time"

	"city-server/internal/utils"
	"city-server/protocol"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	Conn     *websocket.Conn
	Send     chan []byte
	Hub      *Hub
	PlayerID string
}

func (c *Client) ReadPump() {
	defer func() {
		log.Println("[SERVER] 🔌 Client disconnected")
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		log.Println("[SERVER] 🔄 Received pong, extending deadline")
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	log.Println("[SERVER] 📡 ReadPump started")

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("[SERVER] ❌ Read error:", err)
			break
		}

		log.Printf("[SERVER] 📩 Received message: %s", string(message))

		var msg protocol.Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("[SERVER] ⚠️ Invalid JSON:", err)
			continue
		}

		switch msg.Type {
		case "join":
			var join protocol.JoinMessage
			if err := json.Unmarshal(msg.Data, &join); err != nil {
				log.Println("[SERVER] ❌ Failed to parse join:", err)
				continue
			}

			log.Printf("[SERVER] ✅ Player %s joined", join.PlayerID)
			// 1) Добавляем игрока в состояние
			c.PlayerID = join.PlayerID
			c.Hub.Players[join.PlayerID] = protocol.Vec3{X: 0, Y: 0, Z: 0}

			// 2) Отправляем этому клиенту snapshot всех игроков
			var list []protocol.JoinMessage
			for id := range c.Hub.Players {
				// не включаем себя, если не нужно
				if id == join.PlayerID {
					continue
				}
				list = append(list, protocol.JoinMessage{PlayerID: id})
			}
			snapshot := protocol.Message{
				Type: "world_snapshot",
				Data: utils.MustMarshal(protocol.WorldSnapshot{Players: list}),
			}
			c.Send <- utils.MustMarshal(snapshot)

			// 3) Оповещаем остальных о новом игроке
			joined := protocol.Message{
				Type: "player_joined",
				Data: utils.MustMarshal(protocol.PlayerJoined{
					PlayerID: join.PlayerID,
					Position: protocol.Vec3{X: 0, Y: 0, Z: 0},
				}),
			}
			b := utils.MustMarshal(joined)
			for client := range c.Hub.Clients {
				if client != c {
					client.Send <- b
				}
			}

		case "move":
			var move protocol.MoveMessage
			if err := json.Unmarshal(msg.Data, &move); err != nil {
				log.Println("[SERVER] ❌ Failed to parse move message:", err)
				continue
			}

			log.Printf("[SERVER] 🕹️ Player %s moved to: (%.2f, %.2f, %.2f)", move.PlayerID, move.Position.X, move.Position.Y, move.Position.Z)

			c.Hub.Players[move.PlayerID] = move.Position

			playerMoved := protocol.PlayerMoved{
				PlayerID: move.PlayerID,
				Position: move.Position,
			}

			response := protocol.Message{
				Type: "player_moved",
				Data: utils.MustMarshal(playerMoved),
			}
			respBytes := utils.MustMarshal(response)

			log.Printf("[SERVER] 📢 Broadcasting move of %s to all clients", move.PlayerID)
			c.Hub.Broadcast <- respBytes

		default:
			log.Printf("[SERVER] ⚠️ Unknown message type: %s", msg.Type)
		}
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
		log.Println("[SERVER] 📴 WritePump closed")
	}()

	log.Println("[SERVER] 📨 WritePump started")

	for {
		select {
		case msg, ok := <-c.Send:
			if !ok {
				log.Println("[SERVER] ⚠️ Send channel closed")
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("[SERVER] 📤 Sending message: %s", string(msg))
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("[SERVER] ❌ Write error:", err)
				return
			}

		case <-ticker.C:
			log.Println("[SERVER] 🫧 Sending ping")
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("[SERVER] ❌ Ping error:", err)
				return
			}
		}
	}
}
