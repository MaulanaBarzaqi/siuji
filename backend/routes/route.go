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
	pc *controllers.PeriodController,
	sc *controllers.SectionController,
) {
	auth := app.Group("/api/v1/auth")
	auth.Post("/login", ac.Login)
	auth.Post("/forgot-password", ac.ForgotPassword)
	auth.Post("/verify-otp", middleware.TempAuth(utils.PurposeVerifyEmail), ac.VerifyOTP)
	auth.Post("/reset-password", middleware.TempAuth(utils.PurposeResetPassword), ac.ResetPassword)
	auth.Post("/change-password", middleware.JWTAuth(), ac.ChangePassword)
	auth.Get("/me", middleware.JWTAuth(), ac.GetMe)
	auth.Post("/logout", middleware.JWTAuth(), ac.Logout)

	periods := app.Group("/api/v1/periods", middleware.JWTAuth(), middleware.RequireRole("admin"))
	periods.Post("/", pc.CreatePeriod)
	periods.Get("/", pc.GetAllPeriods)
	periods.Get("/:period_public_id", pc.GetPeriodByPublicID)
	periods.Put("/:period_public_id", pc.UpdatePeriod)
	periods.Delete("/:period_public_id", pc.DeletePeriod)
	periods.Post("/:period_public_id/sections", pc.AddSectionToPeriod)
	periods.Delete("/:period_public_id/sections/:section_public_id", pc.RemoveSectionFromPeriod)
	periods.Put("/:period_public_id/questions/reorder", pc.ReorderSections)
	
	sections := app.Group("/api/v1/sections", middleware.JWTAuth(), middleware.RequireRole("admin"))
	sections.Post("/", sc.CreateSection)
	sections.Get("/", sc.GetAllSections)
	sections.Get("/:section_public_id", sc.GetSectionByPublicID)
	sections.Put("/:section_public_id", sc.UpdateSection)
	sections.Delete("/:section_public_id", sc.DeleteSection)
	sections.Put("/:section_public_id/questions/reorder", sc.ReorderQuestions)

	questions := app.Group("/api/v1/questions", middleware.JWTAuth(), middleware.RequireRole("admin"))
	questions.Put("/:question_public_id/answer-key", sc.UpsertAnswerKey)
	questions.Put("/:question_public_id/options/reorder", sc.ReorderOptions)
}
