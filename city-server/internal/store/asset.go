// internal/store/asset.go
package store

import "gorm.io/gorm"

// Asset описывает один загруженный asset‑bundle.
type Asset struct {
	gorm.Model
	WorldID         string `gorm:"index;not null" json:"world_id"`
	Platform        string `gorm:"index;not null" json:"platform"`
	AssetBundleHash string `gorm:"uniqueIndex;not null" json:"asset_bundle_hash"`
	Path            string `gorm:"not null" json:"path"`
}
