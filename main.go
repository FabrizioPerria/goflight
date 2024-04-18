package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()
	api := app.Group("/api")

	api.Get("/greet", handleDefault)

	app.Listen(":5001")
}

func handleDefault(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]string{"/greet": "yo!"})
}
