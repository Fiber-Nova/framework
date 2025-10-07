package controllers

import "github.com/gofiber/fiber/v2"

func Healthz() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	}
}
