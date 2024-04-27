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
	uri = "mongodb://localhost:27017"
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
	var (
		userStore     = db.NewMongoDbUserStore(client)
		userHandler   = handlers.UserHandler{UserStore: userStore}
		flightStore   = db.NewMongoDbFlightStore(client)
		seatStore     = db.NewMongoDbSeatStore(client, *flightStore)
		flightHandler = handlers.FlightHandler{FlightStore: flightStore, SeatStore: seatStore}

		app   = fiber.New(config)
		apiv1 = app.Group("/api/v1/")
	)

	apiv1.Post("/users", userHandler.HandlePostCreateUserv1)
	apiv1.Delete("/users", userHandler.HandleDeleteAllUsersv1)
	apiv1.Get("/users", userHandler.HandleGetUsersv1)

	apiv1.Delete("/users/:id", userHandler.HandleDeleteUserv1)
	apiv1.Get("/users/:id", userHandler.HandleGetUserv1)
	apiv1.Put("/users/:id", userHandler.HandlePutUserv1)

	apiv1.Get("/flights", flightHandler.HandleGetFlightsv1)
	apiv1.Get("/flights/:id/seats", flightHandler.HandleGetSeatsv1)
	// apiv1.Post("/flights", flightHandler.HandlePostCreateFlightv1)

	app.Listen(*listenAddress)
}
