package api

import (
	"github.com/gofiber/fiber/v2"
)

func HandleGetUsersv1(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]string{"user": "dude"})
}

func HandleGetUserv1(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	return ctx.JSON(map[string]string{"user": id})
}
