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
		userStore        = db.NewMongoDbUserStore(client)
		flightStore      = db.NewMongoDbFlightStore(client)
		seatStore        = db.NewMongoDbSeatStore(client, *flightStore)
		reservationStore = db.NewMongoDbReservationStore(client, *flightStore, *seatStore)

		mainStore = db.Store{User: userStore, Flight: flightStore, Seat: seatStore, Reservation: reservationStore}
	)
	app := handlers.SetupRoutes(mainStore, config)
	app.Listen(*listenAddress)
}
