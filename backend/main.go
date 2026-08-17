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
	userRepo := repositories.NewUserRepository(config.DB)
	otpRepo := repositories.NewOTPRepository(config.DB)
	periodRepo := repositories.NewPeriodRepository(config.DB)
	sectionRepo := repositories.NewSectionRepository(config.DB)
	questionRepo := repositories.NewQuestionRepository(config.DB)
	optionRepo := repositories.NewOptionRepository(config.DB)
	participantRepo := repositories.NewParticipantRepository(config.DB)

	// services
	emailService := services.NewEmailService()
	authService := services.NewAuthService(userRepo, otpRepo, emailService)
	periodService := services.NewPeriodService(periodRepo, sectionRepo)
	sectionService := services.NewSectionService(sectionRepo)
	questionService := services.NewQuestionService(questionRepo, sectionRepo)
	optionService := services.NewOptionService(optionRepo, questionRepo)
	userService := services.NewUserService(userRepo)
	participantService := services.NewParticipantService(participantRepo, userRepo, periodRepo)

	// controllers
	authController := controllers.NewAuthController(authService)
	periodController := controllers.NewPeriodController(periodService)
	sectionController := controllers.NewSectionController(sectionService)
	questionController := controllers.NewQuestionController(questionService)
	optionController := controllers.NewOptionController(optionService)
	userController := controllers.NewUserController(userService)
	participantController := controllers.NewParticipantController(participantService)

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

	routes.SetupRoutes(
		app, 
		authController, 
		periodController, 
		sectionController, 
		questionController, 
		optionController,
		participantController,
		userController,
	)

	port := config.AppConfig.Port
	log.Printf("server is running on port:%s", port)
	log.Fatal(app.Listen(":" + port))
}