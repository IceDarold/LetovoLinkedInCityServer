package ws

import (
	"city-server/internal/utils"
	"city-server/protocol"
	"log"
)

type Hub struct {
	Clients    map[*Client]bool
	Players    map[string]protocol.Vec3
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Players:    make(map[string]protocol.Vec3),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Unregister:
			log.Printf("[SERVER] 🔌 Player %s disconnected", client.PlayerID)

			// Remove from world state
			delete(h.Clients, client)
			delete(h.Players, client.PlayerID)

			// Notify others
			leaveMsg := protocol.Message{
				Type: "player_left",
				Data: utils.MustMarshal(protocol.PlayerLeft{
					PlayerID: client.PlayerID,
				}),
			}

			msgBytes := utils.MustMarshal(leaveMsg)
			for c := range h.Clients {
				c.Send <- msgBytes
			}

		case message := <-h.Broadcast:
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
		}
	}
}
