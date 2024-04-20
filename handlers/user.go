package handlers

import (
	"context"

	"github.com/fabrizioperria/goflight/db"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserStore db.UserStorer
}

func (h *UserHandler) HandleGetUserv1(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	backgroundContext := context.Background()
	user, err := h.UserStore.GetUserById(backgroundContext, id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(user)
}

func (h *UserHandler) HandleGetUsersv1(ctx *fiber.Ctx) error {
	backgroundContext := context.Background()
	users, err := h.UserStore.GetUsers(backgroundContext)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(users)
}
