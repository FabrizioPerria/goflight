package handlers

import (
	"os"
	"time"

	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	store db.Store
}
type AuthResponse struct {
	Token string     `json:"token"`
	User  types.User `json:"user"`
}

func NewAuthHandler(store db.Store) *AuthHandler {
	return &AuthHandler{
		store: store,
	}
}

type UserAuthenticate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleAuthenticate(ctx *fiber.Ctx) error {
	var authParams UserAuthenticate
	if err := ctx.BodyParser(&authParams); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	filter := db.Map{"email": authParams.Email}
	user, err := h.store.User.GetUser(ctx.Context(), filter)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	if ok := user.Authenticate(authParams.Password); !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	authResponse := AuthResponse{
		Token: produceToken(user),
		User:  *user,
	}

	return ctx.Status(fiber.StatusOK).JSON(authResponse)
}

func produceToken(user *types.User) string {
	now := time.Now()
	exp := now.Add(time.Hour * 4).UTC().Unix()
	claims := jwt.MapClaims{
		"sub": user.Id.Hex(),
		"exp": exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	return signedToken
}
