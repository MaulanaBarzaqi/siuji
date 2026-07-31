package middleware

import (
	"siuji-backend/utils"

	"github.com/gofiber/fiber/v3"
)

func RequireRole(roles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return utils.Forbidden(c, "Access denied", "Role not found in token")
		}
		for _, allowedRole := range roles {
			if role == allowedRole {
				return c.Next()
			}
		}
		return utils.Forbidden(c, "Access denied", "You do not have permission to access this resource")
	}
}