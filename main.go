package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddress := flag.String("listen", ":5001", "The address to listen on for HTTP requests.")
	flag.Parse()

	app := fiber.New()
	app.Get("/greet", handleDefault)

	apiv1 := app.Group("/api/v1/")
	apiv1.Get("/user", handleUserv1)

	app.Listen(*listenAddress)
}

func handleDefault(ctx *fiber.Ctx) error {
	return ctx.SendString("yo dude!!")
}

func handleUserv1(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]string{"user": "dude"})
}
