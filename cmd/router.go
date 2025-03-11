package main

import (
	"github.com/DosyaKitarov/notification-service/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRouter настраивает маршруты и возвращает настроенный маршрутизатор Gin.
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Эндпоинт проверки работоспособности
	r.GET("/health", handlers.HealthCheck)

	// CRUD эндпоинты для уведомлений
	r.POST("/notifications", handlers.CreateNotification)
	r.GET("/notifications", handlers.GetNotifications)
	r.GET("/notifications/:id", handlers.GetNotification)
	r.PUT("/notifications/:id", handlers.UpdateNotification)
	r.DELETE("/notifications/:id", handlers.DeleteNotification)

	return r
}
