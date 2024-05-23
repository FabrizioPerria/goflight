package handlers

import (
	"time"

	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReservationHandler struct {
	store db.Store
}

func NewReservationHandler(store db.Store) *ReservationHandler {
	return &ReservationHandler{
		store: store,
	}
}

func (h *ReservationHandler) HandlePostCreateReservationv1(ctx *fiber.Ctx) error {
	flightID := ctx.Params("fid")
	fid, err := primitive.ObjectIDFromHex(flightID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	filter := bson.M{"_id": fid}
	flight, err := h.store.Flight.GetFlight(ctx.Context(), filter)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	if len(flight.Seats) == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No seats available"})
	}

	dateFrom, err := time.Parse(time.RFC3339, flight.DepartureTime)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	dateTo, err := time.Parse(time.RFC3339, flight.ArrivalTime)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if time.Now().After(dateFrom) || time.Now().After(dateTo) {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Flight already departed"})
	}

	seatID := ctx.Params("sid")
	sid, err := primitive.ObjectIDFromHex(seatID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user := ctx.Context().UserValue("user").(*types.User)
	reservation, err := h.store.Reservation.CreateReservation(ctx.Context(), bson.M{"_id": sid}, user.Id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(reservation)
}

func (h *ReservationHandler) HandleGetReservationsv1(ctx *fiber.Ctx) error {
	reservations, err := h.store.Reservation.GetReservations(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	user := ctx.Context().UserValue("user").(*types.User)
	if !user.IsAdmin {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	return ctx.JSON(reservations)
}

func (h *ReservationHandler) authenticateUser(ctx *fiber.Ctx, filter bson.M) error {
	reservation, err := h.store.Reservation.GetReservation(ctx.Context(), filter)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	user := ctx.Context().UserValue("user").(*types.User)
	if reservation.UserId != user.Id {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	return nil
}

func (h *ReservationHandler) HandleGetReservationv1(ctx *fiber.Ctx) error {
	reservationID := ctx.Params("rid")
	rid, err := primitive.ObjectIDFromHex(reservationID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	filter := bson.M{"_id": rid}
	err = h.authenticateUser(ctx, filter)
	if err != nil {
		return err
	}
	reservation, err := h.store.Reservation.GetReservation(ctx.Context(), filter)
	return ctx.JSON(reservation)
}

func (h *ReservationHandler) HandleDeleteReservationv1(ctx *fiber.Ctx) error {
	reservationID := ctx.Params("rid")
	rid, err := primitive.ObjectIDFromHex(reservationID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	filter := bson.M{"_id": rid}
	err = h.authenticateUser(ctx, filter)
	if err != nil {
		return err
	}
	if err = h.store.Reservation.DeleteReservation(ctx.Context(), filter); err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Reservation deleted")
}