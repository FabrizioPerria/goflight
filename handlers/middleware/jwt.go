package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/fabrizioperria/goflight/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

func JWTAuthentication(userStore db.UserStorer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, ok := ctx.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		claims, err := parseJWT(token[0])
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}
		expiration := claims["exp"].(float64)
		expirationTime := time.Unix(int64(expiration), 0)
		if time.Now().UTC().After(expirationTime) {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		uid := claims["sub"].(string)
		oid, err := primitive.ObjectIDFromHex(uid)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		filter := db.Map{"_id": oid}
		user, err := userStore.GetUser(ctx.Context(), filter)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		ctx.Context().SetUserValue("user", user)

		return ctx.Next()
	}
}

func parseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("internal server error: unauthorized")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err == nil && token.Valid {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("internal server error: unauthorized")
		}
		return claims, nil
	}

	return nil, fmt.Errorf("internal server error: unauthorized")
}
