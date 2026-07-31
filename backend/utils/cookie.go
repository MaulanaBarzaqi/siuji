package utils

import (
	"siuji-backend/config"
	"time"

	"github.com/gofiber/fiber/v3"
)

// cookie helper
func SetAuthCookies(c fiber.Ctx, accessToken, refreshToken string) {
	accessDuration, _ := time.ParseDuration(config.AppConfig.JWTExpiresIn)
	refreshDuration, _ := time.ParseDuration(config.AppConfig.RefreshTokenExpires)

	c.Cookie(&fiber.Cookie{
		Name: "access_token",
		Value: accessToken,
		HTTPOnly: true,
		Secure: false,  // true jika production
		SameSite: "Lax",
		Path: "/",
		MaxAge: int(accessDuration.Seconds()),
	})

	c.Cookie(&fiber.Cookie{
		Name: "refresh_token",
		Value: refreshToken,
		HTTPOnly: true,
		Secure: false,
		SameSite: "Lax",
		Path: "/api/v1/auth/refresh",
		MaxAge: int(refreshDuration.Seconds()),
	})
}

func ClearAccessTokenCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name: "access_token",
		Value: "",
		HTTPOnly: true,
		Secure: false,
		SameSite: "Lax",
		Path: "/",
		MaxAge: -1,
	})
}

// clear cookie
func ClearAuthCookies(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name: "access_token",
		Value: "",
		HTTPOnly: true,
		MaxAge: -1, // Hapus cookie
	})
	c.Cookie(&fiber.Cookie{
		Name: "refresh_token",
		Value: "",
		HTTPOnly: true,
		Path: "/api/v1/auth/refresh",
		MaxAge: -1,
	})
}