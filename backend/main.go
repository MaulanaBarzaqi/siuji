package main

import (
	"log"
	"siuji-backend/config"
	"siuji-backend/controllers"
	"siuji-backend/database"
	"siuji-backend/database/seed"
	"siuji-backend/middleware"
	"siuji-backend/repositories"
	"siuji-backend/routes"
	"siuji-backend/services"

	"github.com/gofiber/fiber/v3"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	database.Migrate()
	seed.SeedAdmin()

	// repositories
	userRepo := repositories.NewUserRepository()
	otpRepo := repositories.NewOTPRepository()

	// services
	emailService := services.NewEmailService()
	authService := services.NewAuthService(userRepo, otpRepo, emailService)

	// controllers
	authController := controllers.NewAuthController(authService)

	app := fiber.New()

	app.Use(middleware.CORS())
	app.Use(middleware.Logger())

	app.Get("/api/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "success",
			"message": "siuji backend is running!",
			"db": "connected",
		})
	})

	routes.SetupRoutes(app, authController)

	port := config.AppConfig.Port
	log.Printf("server is running on port:%s", port)
	log.Fatal(app.Listen(":" + port))
}