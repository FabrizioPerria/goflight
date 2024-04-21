package main

import (
	"context"
	"flag"
	"log"

	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/handlers"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	uri            = "mongodb://localhost:27017"
	dbName         = "goflight"
	userCollection = "users"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddress := flag.String("listen", ":5001", "The address to listen on for HTTP requests.")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	userStore := db.NewMongoDbUserStore(client, dbName)

	userHandler := handlers.UserHandler{
		UserStore: userStore,
	}

	app := fiber.New(config)

	apiv1 := app.Group("/api/v1/")
	apiv1.Post("/user", userHandler.HandlePostCreateUserv1)
	apiv1.Delete("/user", userHandler.HandleDeleteAllUsersv1)
	apiv1.Get("/user", userHandler.HandleGetUsersv1)

	apiv1.Delete("/user/:id", userHandler.HandleDeleteUserv1)
	apiv1.Get("/user/:id", userHandler.HandleGetUserv1)
	apiv1.Put("/user/:id", userHandler.HandlePutUserv1)

	apiv1.Post("/flight", userHandler.HandlePostCreateUserv1)

	app.Listen(*listenAddress)
}
