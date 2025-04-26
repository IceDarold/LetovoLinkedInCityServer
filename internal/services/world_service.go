package services

import (
	"city-server/internal/store"
	"errors"
	"log"

	"gorm.io/gorm"
)

// WorldService - сервис для работы с миром
type WorldService struct {
	db *gorm.DB
}

// Новый WorldService
func NewWorldService(db *gorm.DB) *WorldService {
	return &WorldService{
		db: db,
	}
}

// GetWorldState - получает состояние мира для заданного worldId и платформы
func (ws *WorldService) GetWorldState(worldId, platform string) (map[string]interface{}, error) {
	var world store.World
	// Ищем мир в базе данных по идентификатору и платформе
	err := ws.db.Where("world_id = ? AND platform = ?", worldId, platform).First(&world).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("мир не найден")
		}
		log.Printf("Ошибка при получении состояния мира: %s", err)
		return nil, err
	}

	// Преобразуем состояние мира в удобный формат (например, мапу для ответа)
	state := map[string]interface{}{
		"world_id": world.WorldID,
		"platform": world.Platform,
		"state":    world.State, // предполагаем, что поле State содержит JSON или структуру данных
	}

	return state, nil
}

// SaveWorldState - сохраняет состояние мира в базе данных
func (ws *WorldService) SaveWorldState(worldId, platform string, state map[string]interface{}) error {
	var world store.World

	// Проверяем, существует ли уже мир
	err := ws.db.Where("world_id = ? AND platform = ?", worldId, platform).First(&world).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Ошибка при проверке существования мира: %s", err)
		return err
	}

	// Если мир не найден, создаём новый
	if errors.Is(err, gorm.ErrRecordNotFound) {
		world = store.World{
			WorldID:  worldId,
			Platform: platform,
			State:    state, // Преобразуем state в формат, который можно сохранить
		}

		// Сохраняем новый мир
		if err := ws.db.Create(&world).Error; err != nil {
			log.Printf("Ошибка при создании нового мира: %s", err)
			return err
		}
	} else {
		// Если мир найден, обновляем его состояние
		world.State = state
		if err := ws.db.Save(&world).Error; err != nil {
			log.Printf("Ошибка при обновлении мира: %s", err)
			return err
		}
	}

	return nil
}

// GetWorldDelta - получает изменения (патч) мира с момента последнего состояния
func (ws *WorldService) GetWorldDelta(worldId, platform, lastKnownSnapshotHash string) (map[string]interface{}, error) {
	// Для простоты предположим, что изменения (патч) хранятся как разница между состоянием мира
	var world store.World
	err := ws.db.Where("world_id = ? AND platform = ?", worldId, platform).First(&world).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("мир не найден")
		}
		log.Printf("Ошибка при получении патча мира: %s", err)
		return nil, err
	}

	// Например, здесь можно реализовать логику для сравнения текущего состояния и предыдущего,
	// чтобы возвращать только изменения.
	// Для простоты вернем текущее состояние как патч:
	patch := map[string]interface{}{
		"world_id": world.WorldID,
		"platform": world.Platform,
		"delta":    world.State, // Можно сравнить с предыдущим состоянием
	}

	return patch, nil
}
