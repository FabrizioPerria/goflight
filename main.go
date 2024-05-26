package main

import (
	"context"
	"flag"
	"log"

	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/handlers"
	"github.com/fabrizioperria/goflight/handlers/middleware"
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

		userHandler        = handlers.NewUserHandler(mainStore)
		flightHandler      = handlers.NewFlightHandler(mainStore)
		authHandler        = handlers.NewAuthHandler(mainStore)
		reservationHandler = handlers.NewReservationHandler(mainStore)

		app     = fiber.New(config)
		notAuth = app.Group("/api")
		apiv1   = app.Group("/api/v1/", middleware.JWTAuthentication(userStore))
		admin   = apiv1.Group("/admin", middleware.AdminOnly())
	)
	notAuth.Post("/auth", authHandler.HandleAuthenticate)
	notAuth.Post("/v1/users", userHandler.HandlePostCreateUserv1)

	admin.Post("/users", userHandler.HandlePostCreateAdminUserv1)
	admin.Delete("/users", userHandler.HandleDeleteAllUsersv1)
	admin.Get("/users", userHandler.HandleGetUsersv1)
	admin.Post("/flights", flightHandler.HandlePostCreateFlightv1)
	admin.Get("/reservations", reservationHandler.HandleGetAllReservationsv1)

	apiv1.Delete("/users/:uid", userHandler.HandleDeleteUserv1)
	apiv1.Get("/users/:uid", userHandler.HandleGetUserv1)
	apiv1.Put("/users/:uid", userHandler.HandlePutUserv1)

	apiv1.Get("/flights", flightHandler.HandleGetFlightsv1)
	apiv1.Get("/flights/:fid", flightHandler.HandleGetFlightv1)
	apiv1.Get("/flights/:fid/seats", flightHandler.HandleGetSeatsv1)
	apiv1.Get("/flights/:fid/seats/:sid", flightHandler.HandleGetSeatv1)

	apiv1.Post("/flights/:fid/seats/:sid/reservations", reservationHandler.HandlePostCreateReservationv1)

	apiv1.Get("/reservations", reservationHandler.HandleGetMyReservationsv1)

	apiv1.Get("/reservations/:rid", reservationHandler.HandleGetReservationv1)
	apiv1.Delete("/reservations/:rid", reservationHandler.HandleDeleteReservationv1)

	app.Listen(*listenAddress)
}
