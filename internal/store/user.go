// internal/store/user.go
package store

import "gorm.io/gorm"

// User хранит данные учётной записи игрока/разработчика.
type User struct {
	gorm.Model
	UUID         string `gorm:"uniqueIndex;not null" json:"uuid"`
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	Role         string `gorm:"default:'player'" json:"role"`
}
