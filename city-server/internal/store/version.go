// internal/store/version.go
package store

import "gorm.io/gorm"

// Version хранит информацию о выложенных snapshot’ах мира.
type Version struct {
	gorm.Model
	WorldID      string `gorm:"index;not null" json:"world_id"`
	Platform     string `gorm:"index;not null" json:"platform"`
	SnapshotHash string `gorm:"uniqueIndex:idx_version_world_platform;not null" json:"snapshot_hash"`
}
