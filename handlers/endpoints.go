package handlers

import (
	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/handlers/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(mainStore db.Store, config fiber.Config) *fiber.App {
	userHandler := NewUserHandler(mainStore)
	flightHandler := NewFlightHandler(mainStore)
	authHandler := NewAuthHandler(mainStore)
	reservationHandler := NewReservationHandler(mainStore)

	app := fiber.New(config)
	notAuth := app.Group("/api")
	apiv1 := app.Group("/api/v1/", middleware.JWTAuthentication(mainStore.User))
	admin := apiv1.Group("/admin", middleware.AdminOnly())

	notAuth.Post("/auth", authHandler.HandleAuthenticate)
	notAuth.Post("/users", userHandler.HandlePostCreateUserv1)

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

	return app
}
