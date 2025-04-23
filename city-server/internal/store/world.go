// internal/store/world.go
package store

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// World хранит текущее состояние игрового мира.
type World struct {
	gorm.Model
	WorldID  string            `gorm:"uniqueIndex:idx_world_platform;not null" json:"world_id"`
	Platform string            `gorm:"uniqueIndex:idx_world_platform;not null" json:"platform"`
	State    datatypes.JSONMap `gorm:"type:jsonb"`
}
