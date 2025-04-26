package store

import "time"

type AuthToken struct {
	Token     string `gorm:"primaryKey"`
	PlayerID  string
	CreatedAt time.Time
}
