package middleware

import (
	"log"
	"siuji-backend/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func extractTempToken(c fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}
	return ""
}

func TempAuth(expectedPurpose string) fiber.Handler {
	return func(c fiber.Ctx) error {
		tokenString := extractTempToken(c)
		if tokenString == "" {
			log.Printf("[TEMP_AUTH] Missing temp token - IP: %s, Path: %s", c.IP(), c.Path())
			return utils.Unauthorized(c, "Unauthorized", "Missing temp token in Authorization header")
		}

		claims, err := utils.ValidateTempToken(tokenString, expectedPurpose)
		if err != nil {
			log.Printf("[TEMP_AUTH] Invalid temp token - IP: %s, Path: %s, Error: %v", 
				c.IP(), c.Path(), err)
			return utils.Unauthorized(c, "Unauthorized", "Invalid or expired temp token")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("token_purpose", claims.Purpose)

		log.Printf("[TEMP_AUTH] Success - Purpose: %s, UserID: %d, Email: %s, IP: %s", 
			claims.Purpose, claims.UserID, claims.Email, c.IP())
		return c.Next()
	}
}