package api

import (
	"github.com/fabrizioperria/goflight/types"
	"github.com/gofiber/fiber/v2"
)

func HandleGetUsersv1(ctx *fiber.Ctx) error {
	u := types.User{
		FirstName: "John",
		LastName:  "Doe",
	}
	return ctx.JSON(u)
}

func HandleGetUserv1(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	return ctx.JSON(map[string]string{"user": id})
}
