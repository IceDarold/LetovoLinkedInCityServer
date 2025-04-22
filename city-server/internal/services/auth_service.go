package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"

	"github.com/spf13/viper"
)

// AuthService отвечает за валидацию HMAC‑подписей от клиента
type AuthService struct {
	hmacSecret []byte
}

// NewAuthService подгружает секрет из конфига и возвращает AuthService
func NewAuthService() *AuthService {
	secret := viper.GetString("auth.hmac_secret")
	if secret == "" {
		log.Fatal("AuthService: missing config value 'auth.hmac_secret'")
	}
	return &AuthService{
		hmacSecret: []byte(secret),
	}
}

// ValidateSignature проверяет, что переданная hex‑строка signature соответствует
// HMAC‑SHA256(payload) с использованием секрета из конфига.
func (s *AuthService) ValidateSignature(signature, payload string) (bool, error) {
	// Декодируем переданную подпись из hex
	sigBytes, err := hex.DecodeString(signature)
	if err != nil {
		return false, errors.New("invalid signature format")
	}

	// Вычисляем HMAC‑SHA256 от payload
	mac := hmac.New(sha256.New, s.hmacSecret)
	if _, err := mac.Write([]byte(payload)); err != nil {
		return false, err
	}
	expectedMAC := mac.Sum(nil)

	// Сравниваем безопасно
	if !hmac.Equal(expectedMAC, sigBytes) {
		return false, nil
	}
	return true, nil
}
