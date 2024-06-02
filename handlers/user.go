package handlers

import (
	"fmt"

	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	store db.Store
}

func NewUserHandler(store db.Store) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (h *UserHandler) HandleGetUserv1(ctx *fiber.Ctx) error {
	id := ctx.Params("uid")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	filter := db.Map{"_id": oid}

	user, err := h.store.User.GetUser(ctx.Context(), filter)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(user)
}

func (h *UserHandler) HandleGetUsersv1(ctx *fiber.Ctx) error {
	pagination := db.Pagination{
		Page:  ctx.Query("page"),
		Limit: ctx.Query("limit"),
	}
	users, err := h.store.User.GetUsers(ctx.Context(), &pagination)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(users)
}

func (h *UserHandler) HandlePostCreateUserv1(ctx *fiber.Ctx) error {
	fmt.Println("HandlePostCreateUserv1")
	createUserParams := types.CreateUserParams{}
	err := ctx.BodyParser(&createUserParams)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	errors := createUserParams.Validate()
	if len(errors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors})
	}

	user, err := types.NewUserFromParams(createUserParams, false)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = h.store.User.CreateUser(ctx.Context(), user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) HandlePostCreateAdminUserv1(ctx *fiber.Ctx) error {
	createUserParams := types.CreateUserParams{}
	err := ctx.BodyParser(&createUserParams)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	errors := createUserParams.Validate()
	if len(errors) > 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errors})
	}

	user, err := types.NewUserFromParams(createUserParams, true)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = h.store.User.CreateUser(ctx.Context(), user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusCreated).JSON(user)
}

func (h *UserHandler) HandleDeleteUserv1(ctx *fiber.Ctx) error {
	userId := ctx.Params("uid")
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	filter := db.Map{"_id": oid}

	id, err := h.store.User.DeleteUser(ctx.Context(), filter)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("User deleted: " + id)
}

func (h *UserHandler) HandleDeleteAllUsersv1(ctx *fiber.Ctx) error {
	err := h.store.User.Drop(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Users deleted")
}

func (h *UserHandler) HandlePutUserv1(ctx *fiber.Ctx) error {
	userID := ctx.Params("uid")
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	filter := db.Map{"_id": oid}

	values := types.UpdateUserParams{}
	err = ctx.BodyParser(&values)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	_, err = h.store.User.UpdateUser(ctx.Context(), filter, values)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("User updated: " + userID)
}
