package main

import (
	"flag"

	api "github.com/fabrizioperria/goflight/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddress := flag.String("listen", ":5001", "The address to listen on for HTTP requests.")
	flag.Parse()

	app := fiber.New()
	app.Get("/greet", handleDefault)

	apiv1 := app.Group("/api/v1/")
	apiv1.Get("/user", api.HandleGetUsersv1)
	apiv1.Get("/user/:id", api.HandleGetUserv1)

	app.Listen(*listenAddress)
}

func handleDefault(ctx *fiber.Ctx) error {
	return ctx.SendString("yo dude!")
}
