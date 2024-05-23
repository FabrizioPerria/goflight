package middleware

import (
	"github.com/fabrizioperria/goflight/types"
	"github.com/gofiber/fiber/v2"
)

func AdminOnly() func(*fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		user := ctx.Context().UserValue("user").(*types.User)
		if !user.IsAdmin {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}
		return ctx.Next()
	}
}
