package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"city-server/internal/api"
	"city-server/internal/middleware"
	"city-server/internal/services"
	"city-server/internal/store"
	"city-server/internal/ws"
)

var db *gorm.DB
var err error

func init() {
	// читаем конфиг
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Ошибка чтения конфига: %s", err)
	}

	// открываем GORM v2
	dsn := viper.GetString("database.dsn")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %s", err)
	}

	// мигрируем модели
	if err := db.AutoMigrate(&store.World{}, &store.User{}, &store.Asset{}, &store.Version{}); err != nil {
		log.Fatalf("Ошибка миграции БД: %s", err)
	}

	log.Println("✅ Подключение к PostgreSQL успешно")
}

func main() {
	// Инициализация сервисов
	worldService := services.NewWorldService(db)
	assetService := services.NewAssetService(db)
	authService := services.NewAuthService()
	notificationService := services.NewNotificationService()

	// Инициализация WebSocket-хаба
	wsHub := ws.NewHub()
	go wsHub.Run()

	// Главный роутер
	r := mux.NewRouter()

	// 💬 WebSocket маршрут без middleware
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWS(wsHub, w, r)
	})

	// 📦 Подроутер для API с middleware
	apiRouter := r.PathPrefix("/").Subrouter()
	apiRouter.Use(middleware.LoggingMiddleware)
	apiRouter.Use(middleware.ErrorMiddleware(notificationService))
	token := viper.GetString("auth.api_token")
	apiRouter.Use(middleware.AuthMiddleware(token))

	// Инициализация API обработчиков
	apiHandler := api.NewHandler(worldService, assetService, authService)

	// API маршруты
	apiRouter.HandleFunc("/world/{worldId}/state/{platform}", apiHandler.GetWorldState).Methods(http.MethodGet)
	apiRouter.HandleFunc("/world/{worldId}/state/{platform}", apiHandler.SaveWorldState).Methods(http.MethodPost)
	apiRouter.HandleFunc("/world/{worldId}/delta/{platform}/{lastKnownSnapshotHash}", apiHandler.GetWorldDelta).Methods(http.MethodGet)
	apiRouter.HandleFunc("/assets/{assetBundleHash}", apiHandler.GetAssetBundle).Methods(http.MethodGet)
	apiRouter.HandleFunc("/assets/upload/{worldId}/{platform}/{assetBundleHash}", apiHandler.UploadAssetBundle).Methods(http.MethodPost)
	apiRouter.HandleFunc("/auth/validate-signature/{platform}", apiHandler.ValidateSignature).Methods(http.MethodPost)

	// HTTP-сервер
	server := &http.Server{
		Addr:           viper.GetString("server.address"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("🚀 Запуск сервера на %s...\n", viper.GetString("server.address"))
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("❌ Ошибка запуска сервера: %s", err)
	}
}
