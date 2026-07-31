package middleware

import (
	"github.com/gofiber/fiber/v3"
	fiberlogger "github.com/gofiber/fiber/v3/middleware/logger"
)

func Logger() fiber.Handler {
	return fiberlogger.New(fiberlogger.Config{
		Format: "${time} | ${status} | ${latency} | ${method} ${path} | ${ip}\n",
	})
}