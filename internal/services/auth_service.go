package services

import (
	"city-server/internal/store"
	"log"

	"gorm.io/gorm"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) ValidateToken(token string) (string, bool) {
	var auth store.AuthToken
	result := s.db.First(&auth, "token = ?", token)

	if result.Error != nil {
		log.Printf("[AUTH] ❌ Token not found: %s", token)
		return "", false
	}

	return auth.PlayerID, true
}
