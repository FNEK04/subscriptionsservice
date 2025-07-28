package main

import (
	"database/sql"
	"subscriptionsservice/docs" // Импорт сгенерированной Swagger-документации
	"subscriptionsservice/internal/config"
	"subscriptionsservice/internal/handlers"
	"subscriptionsservice/internal/repositories"
	"subscriptionsservice/internal/services"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"     // Правильный импорт для swaggerFiles
	ginSwagger "github.com/swaggo/gin-swagger" // Пакет для интеграции Swagger с Gin
)

// @title Subscriptions Service API
// @version 1.0
// @description API для управления подписками
// @host localhost:8080
// @BasePath /

var log = logrus.New()

func init() {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	connStr := "user=" + cfg.DBUser + " dbname=" + cfg.DBName + " password=" + cfg.DBPassword + " host=" + cfg.DBHost + " port=" + cfg.DBPort
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	defer db.Close()

	repo := repositories.NewPostgresSubscriptionRepository(db)
	service := services.NewSubscriptionService(repo)
	handler := handlers.NewHandler(service)

	r := gin.Default()
	r.POST("/subscriptions", handler.CreateSubscription)
	r.GET("/subscriptions/:id", handler.GetSubscription)
	r.PUT("/subscriptions/:id", handler.UpdateSubscription)
	r.DELETE("/subscriptions/:id", handler.DeleteSubscription)
	r.GET("/subscriptions", handler.ListSubscriptions)
	r.POST("/subscriptions/total-cost", handler.CalculateTotalCost)

	// Swagger
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Info("Starting server on port ", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}
