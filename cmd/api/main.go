// @title PR Reviewer Assignment Service
// @version 1.0.0
// @description Микросервис для автоматического назначения ревьюеров на Pull Request'ы

// @contact.name API Support
// @contact.url http://localhost:8080
// @contact.email shkeera@mail.ru

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"ReviewAssigner/internal/config"
	"ReviewAssigner/internal/database"
	"ReviewAssigner/internal/handler"
	"ReviewAssigner/internal/repository"
	"ReviewAssigner/internal/service"
	"ReviewAssigner/logger"

	_ "ReviewAssigner/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Summary Health check
// @Description Проверка работоспособности сервиса
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "status: ok"
// @Router /health [get]
func main() {
	cfg := config.Load()

	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// репозитории
	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	prRepo := repository.NewPRRepository(db)

	logger.Init("development") // или "production"
	// сервисы
	reviewService := service.NewReviewService(userRepo, prRepo, logger.Logger)
	prService := service.NewPRService(prRepo, userRepo, reviewService, logger.Logger)
	userService := service.NewUserService(userRepo, teamRepo, prRepo, reviewService, logger.Logger)
	teamService := service.NewTeamService(teamRepo, userRepo, logger.Logger)

	handlers := handler.NewHandler(teamService, userService, prService)

	router := gin.Default()

	router.Use(gin.Logger())

	// тз
	router.StaticFile("/specs/technical-task.yaml", "./openapi.yml")
	router.GET("/docs/technical-task/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/specs/technical-task.yaml"),
		ginSwagger.InstanceName("technical-task"),
	))

	// реализация
	router.GET("/docs/api/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	handlers.SetupRoutes(router)

	// сервер
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := router.Run(":" + cfg.ServerPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //Ctrl+C , kill / docker-compose stop
	<-quit

	log.Println("Shutting down server...")
	db.Close()
	log.Println("Server stopped")
}
