package services

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	"city-server/internal/store"
)

// AssetService отвечает за сохранение и выдачу asset‑bundle’ов
type AssetService struct {
	db          *gorm.DB
	storagePath string
}

// NewAssetService создаёт AssetService, подгружая путь для хранения из конфигурации
func NewAssetService(db *gorm.DB) *AssetService {
	// В config.yaml:
	// assets:
	//   storage_path: "./data/assets"
	storagePath := viper.GetString("assets.storage_path")
	if storagePath == "" {
		log.Fatal("Не задан путь assets.storage_path в конфиге")
	}
	return &AssetService{
		db:          db,
		storagePath: storagePath,
	}
}

// SaveAssetBundle сохраняет бинарник bundle в файловую систему и метаданные в БД
func (as *AssetService) SaveAssetBundle(worldID, platform, hash string, data []byte) error {
	// 1. Подготовить директорию: <storagePath>/<worldID>/<platform>/
	dir := filepath.Join(as.storagePath, worldID, platform)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Printf("Не удалось создать директорию '%s': %v", dir, err)
		return err
	}

	// 2. Записать файл: <hash>.bundle
	filename := fmt.Sprintf("%s.bundle", hash)
	fullPath := filepath.Join(dir, filename)
	if err := ioutil.WriteFile(fullPath, data, 0644); err != nil {
		log.Printf("Не удалось записать файл '%s': %v", fullPath, err)
		return err
	}

	// 3. Сохранить/обновить метаданные в БД
	var asset store.Asset
	err := as.db.Where("asset_bundle_hash = ?", hash).First(&asset).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Ошибка при поиске метаданных asset '%s': %v", hash, err)
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// создаём новую запись
		asset = store.Asset{
			WorldID:         worldID,
			Platform:        platform,
			AssetBundleHash: hash,
			Path:            fullPath,
		}
		if err := as.db.Create(&asset).Error; err != nil {
			log.Printf("Не удалось создать метаданные asset: %v", err)
			return err
		}
	} else {
		// обновляем существующую запись (возможно путь изменился)
		asset.Path = fullPath
		if err := as.db.Save(&asset).Error; err != nil {
			log.Printf("Не удалось обновить метаданные asset: %v", err)
			return err
		}
	}

	return nil
}

// GetAssetBundle возвращает бинарные данные bundle по его хэшу
func (as *AssetService) GetAssetBundle(hash string) ([]byte, error) {
	var asset store.Asset
	if err := as.db.Where("asset_bundle_hash = ?", hash).First(&asset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("asset '%s' не найден", hash)
		}
		log.Printf("Ошибка при запросе метаданных asset '%s': %v", hash, err)
		return nil, err
	}

	data, err := ioutil.ReadFile(asset.Path)
	if err != nil {
		log.Printf("Не удалось прочитать файл '%s': %v", asset.Path, err)
		return nil, err
	}
	return data, nil
}
