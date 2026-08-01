package main

import (
	"log"
	"siuji-backend/config"
	"siuji-backend/controllers"
	"siuji-backend/database"
	"siuji-backend/database/seed"
	_ "siuji-backend/docs"
	"siuji-backend/middleware"
	"siuji-backend/repositories"
	"siuji-backend/routes"
	"siuji-backend/services"

	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
)

// @title           SIUJI API Documentation
// @version         1.0
// @description     API Backend untuk aplikasi SIUJI
// @termsOfService  http://swagger.io/terms/

// @contact.name    Support Team
// @contact.email   support@example.com

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /
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

	app.Get("/swagger/*", swaggo.HandlerDefault)
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