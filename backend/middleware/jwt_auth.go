package middleware

import (
	"log"
	"siuji-backend/utils"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func extractToken(c fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}
	return c.Cookies("access_token")
}

func JWTAuth() fiber.Handler {
	return func (c fiber.Ctx) error {
		tokenString := extractToken(c)
		if tokenString == "" {
			log.Printf("[AUTH] Missing token - IP: %s, Path: %s", c.IP(), c.Path())
			return utils.Unauthorized(c, "Unauthorized", "Missing token")
		}

		claims, err := utils.ValidateAccessToken(tokenString)
		if err != nil {
			log.Printf("[AUTH] Invalid token - IP: %s, Error: %v", c.IP(), err)
			utils.ClearAccessTokenCookie(c)
			return utils.Unauthorized(c, "Unauthorized", "invalid or expired token")
		}
		log.Printf("[AUTH] Success - User: %d, IP: %s", claims.UserID, c.IP())
		c.Locals("user_id", claims.UserID)
		c.Locals("public_id", claims.PublicID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)

		return c.Next()
	}
	
}