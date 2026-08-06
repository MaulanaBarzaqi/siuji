package routes

import (
	"siuji-backend/controllers"
	"siuji-backend/middleware"
	"siuji-backend/utils"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(
	app *fiber.App,
	ac *controllers.AuthController,
	uc *controllers.UserController,
) {
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", ac.Login)
	auth.Post("/forgot-password", ac.ForgotPassword)
	auth.Post("/verify-otp", middleware.TempAuth(utils.PurposeVerifyEmail), ac.VerifyOTP)
	auth.Post("/reset-password", middleware.TempAuth(utils.PurposeResetPassword), ac.ResetPassword)
	auth.Post("/change-password", middleware.JWTAuth(), ac.ChangePassword)
	auth.Get("/me", middleware.JWTAuth(), ac.GetMe)
	auth.Post("/logout", middleware.JWTAuth(), ac.Logout)

	participant := app.Group("/api/v1/participants", middleware.JWTAuth(), middleware.RequireRole("admin"))
	participant.Post("/", uc.CreateParticipant)
	participant.Post("/import", uc.ImportParticipants)
	participant.Get("/", uc.GetAllParticipants)
	participant.Put("/:id", uc.UpdateParticipant)
	participant.Delete("/:id", uc.DeleteParticipant)
}
