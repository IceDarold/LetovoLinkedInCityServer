package services

import (
	"city-server/internal/ws"
)

type StatsService struct {
	hub *ws.Hub
}

func NewStatsService(hub *ws.Hub) *StatsService {
	return &StatsService{hub: hub}
}

func (s *StatsService) GetServerStatus() map[string]interface{} {
	players := []map[string]interface{}{}

	for playerID, pos := range s.hub.Players {
		if playerID != "" {
			players = append(players, map[string]interface{}{
				"playerId": playerID,
				"position": pos,
			})
		}

	}

	return map[string]interface{}{
		"count":   len(s.hub.Clients),
		"players": players,
	}
}
