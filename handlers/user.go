package handlers

import (
	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserStore db.UserStorer
}

func (h *UserHandler) HandleGetUserv1(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := h.UserStore.GetUserById(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(user)
}

func (h *UserHandler) HandleGetUsersv1(ctx *fiber.Ctx) error {
	users, err := h.UserStore.GetUsers(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(users)
}

func (h *UserHandler) HandleCreateRandomUserv1(ctx *fiber.Ctx) error {
	user, err := h.UserStore.CreateRandomUser(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) HandlePostCreateUserv1(ctx *fiber.Ctx) error {
	createUserParams := types.CreateUserParams{}
	err := ctx.BodyParser(&createUserParams)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := types.NewUserFromParams(createUserParams)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	_, err = h.UserStore.CreateUser(ctx.Context(), user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) HandleDeleteAllUsersv1(ctx *fiber.Ctx) error {
	err := h.UserStore.DeleteUsers(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Users deleted")
}
