package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"

	"city-server/internal/api"
	"city-server/internal/middleware"
	"city-server/internal/services"
	"city-server/internal/store"
)

var db *gorm.DB
var err error

func init() {
	// Загружаем конфигурацию из файла config.yaml
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Ошибка чтения конфигурации: %s", err)
	}

	// Подключаемся к базе данных Postgres
	db, err = gorm.Open("postgres", viper.GetString("database.dsn"))
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %s", err)
	}

	// Автоматически мигрируем структуру базы данных
	err = db.AutoMigrate(&store.World{}, &store.User{}, &store.Asset{}, &store.Version{}).Error
	if err != nil {
		log.Fatalf("Ошибка миграции базы данных: %s", err)
	}

	// Логируем успешное подключение
	log.Println("Подключение к базе данных PostgreSQL установлено успешно")
}

func main() {

	// Инициализация сервисов
	worldService := services.NewWorldService(db)
	assetService := services.NewAssetService(db)
	authService := services.NewAuthService()
	notificationService := services.NewNotificationService()

	// Настроим маршруты
	r := mux.NewRouter()
	// передаём notificationService в мидлварь
	r.Use(middleware.ErrorMiddleware(notificationService))

	// Инициализация обработчиков
	apiHandler := api.NewHandler(worldService, assetService, authService)

	// Настроим middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.ErrorMiddleware)

	// API маршруты
	r.HandleFunc("/world/{worldId}/state/{platform}", apiHandler.GetWorldState).Methods(http.MethodGet)
	r.HandleFunc("/world/{worldId}/state/{platform}", apiHandler.SaveWorldState).Methods(http.MethodPost)
	r.HandleFunc("/world/{worldId}/delta/{platform}/{lastKnownSnapshotHash}", apiHandler.GetWorldDelta).Methods(http.MethodGet)
	r.HandleFunc("/assets/{assetBundleHash}", apiHandler.GetAssetBundle).Methods(http.MethodGet)
	r.HandleFunc("/assets/upload/{worldId}/{platform}/{assetBundleHash}", apiHandler.UploadAssetBundle).Methods(http.MethodPost)
	r.HandleFunc("/auth/validate-signature/{platform}", apiHandler.ValidateSignature).Methods(http.MethodPost)

	// Стартуем HTTP сервер
	server := &http.Server{
		Addr:           viper.GetString("server.address"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	log.Printf("Запуск сервера на %s...\n", viper.GetString("server.address"))
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %s", err)
	}
}
