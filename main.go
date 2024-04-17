package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Get("/", handleDefault)

	app.Listen(":5001")
}

func handleDefault(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]string{"/": "default endpoint"})
}
