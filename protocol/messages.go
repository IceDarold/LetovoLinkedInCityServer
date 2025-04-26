package protocol

import "encoding/json"

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// Позиция игрока
type Vec3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

// Позиция игрока
type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Пришёл клиент: JOIN
type JoinMessage struct {
	PlayerID string `json:"playerId"`
	Position Vec3   `json:"playerPosition"`
}

// MOVE от клиента
type MoveMessage struct {
	PlayerID string `json:"playerId"`
	Position Vec3   `json:"position"`
}

// Уведомление другим клиентам
type PlayerMoved struct {
	PlayerID string `json:"playerId"`
	Position Vec3   `json:"position"`
}

// snapshot (при подключении)
type WorldSnapshot struct {
	Players []JoinMessage `json:"players"`
}

// уведомление о новом игроке
type PlayerJoined struct {
	PlayerID string `json:"playerId"`
	Position Vec3   `json:"position"`
}

type PlayerLeft struct {
	PlayerID string `json:"playerId"`
}

type InputMessage struct {
	PlayerID string `json:"playerId"`
	Move     Vec2   `json:"move"`
	Jump     bool   `json:"jump"`
	Sprint   bool   `json:"sprint"`
	Dance    bool   `json:"dance"`
	Look     Vec2   `json:"look"`
}
