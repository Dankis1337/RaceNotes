package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/levferril/racenotes/backend/internal/handlers"
	"github.com/levferril/racenotes/backend/internal/middleware"
	"github.com/levferril/racenotes/backend/internal/models"
	"github.com/levferril/racenotes/backend/internal/services"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "racenotes")
	dbPassword := getEnv("DB_PASSWORD", "racenotes")
	dbName := getEnv("DB_NAME", "racenotes")
	jwtSecret := getEnv("JWT_SECRET", "racenotes-secret-key")
	port := getEnv("PORT", "8080")

	middleware.JWTSecret = []byte(jwtSecret)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&models.User{}, &models.Setup{}, &models.Race{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	authService := services.NewAuthService(db)
	authHandler := handlers.NewAuthHandler(authService)
	setupService := services.NewSetupService(db)
	setupHandler := handlers.NewSetupHandler(setupService)
	raceService := services.NewRaceService(db)
	raceHandler := handlers.NewRaceHandler(raceService)
	calcHandler := handlers.NewCalculatorHandler()
	uploadHandler := handlers.NewUploadHandler()

	// Start email reminder cron
	notificationService := services.NewNotificationService(db)
	notificationService.StartReminderCron()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	user := api.Group("/user")
	user.Use(middleware.AuthRequired())
	{
		user.GET("/profile", authHandler.GetProfile)
		user.PUT("/profile", authHandler.UpdateProfile)
	}

	setups := api.Group("/setups")
	setups.Use(middleware.AuthRequired())
	{
		setups.POST("", setupHandler.Create)
		setups.GET("", setupHandler.List)
		setups.GET("/:id", setupHandler.GetByID)
		setups.PUT("/:id", setupHandler.Update)
		setups.DELETE("/:id", setupHandler.Delete)
	}

	races := api.Group("/races")
	races.Use(middleware.AuthRequired())
	{
		races.POST("", raceHandler.Create)
		races.GET("", raceHandler.List)
		races.GET("/:id", raceHandler.GetByID)
		races.PUT("/:id", raceHandler.Update)
		races.DELETE("/:id", raceHandler.Delete)
	}

	calculator := api.Group("/calculator")
	calculator.Use(middleware.AuthRequired())
	{
		calculator.POST("/tire-pressure", calcHandler.Calculate)
	}

	upload := api.Group("/upload")
	upload.Use(middleware.AuthRequired())
	{
		upload.POST("", uploadHandler.Upload)
	}

	// Serve uploaded files
	api.Static("/uploads", "./uploads")

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
