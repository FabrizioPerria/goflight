package handlers

import (
	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlightHandler struct {
	store db.Store
}

func NewFlightHandler(store db.Store) *FlightHandler {
	return &FlightHandler{
		store: store,
	}
}

func (h *FlightHandler) HandleGetFlightv1(ctx *fiber.Ctx) error {
	id := ctx.Params("fid")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	filter := bson.M{"_id": oid}
	flight, err := h.store.Flight.GetFlight(ctx.Context(), filter)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(flight)
}

func (h *FlightHandler) HandleGetSeatsv1(ctx *fiber.Ctx) error {
	flightID := ctx.Params("fid")
	oid, err := primitive.ObjectIDFromHex(flightID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	filter := bson.M{"flight_id": oid, "available": true}
	seats, err := h.store.Seat.GetSeats(ctx.Context(), filter)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(seats)
}

func (h *FlightHandler) HandleGetFlightsv1(ctx *fiber.Ctx) error {
	flights, err := h.store.Flight.GetFlights(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(flights)
}

func (h *FlightHandler) HandlePostCreateFlightv1(ctx *fiber.Ctx) error {
	createFlightParams := types.CreateFlightParams{}
	err := ctx.BodyParser(&createFlightParams)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	flight, err := types.NewFlightFromParams(createFlightParams)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = h.store.Flight.CreateFlight(ctx.Context(), flight)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	seatIDs := []primitive.ObjectID{}
	for i := 0; i < createFlightParams.NumberOfSeats; i++ {
		seat := types.Seat{
			FlightId:  primitive.ObjectID(flight.Id),
			Number:    i,
			Price:     100,
			Class:     types.SeatClass(i%3 + 1),
			Location:  types.SeatLocation(i%3 + 1),
			Available: true,
		}
		created, err := h.store.Seat.CreateSeat(ctx.Context(), &seat)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		seatIDs = append(seatIDs, created.Id)
	}
	h.store.Flight.UpdateFlight(ctx.Context(), bson.M{"_id": flight.Id}, types.UpdateFlightParams{Seats: seatIDs})
	flight.Seats = seatIDs

	return ctx.Status(fiber.StatusCreated).JSON(flight)
}

func (h *FlightHandler) HandlePutFlightv1(ctx *fiber.Ctx) error {
	flightID := ctx.Params("fid")
	oid, err := primitive.ObjectIDFromHex(flightID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	filter := bson.M{"_id": oid}
	updateFlightParams := types.UpdateFlightParams{}

	err = ctx.BodyParser(&updateFlightParams)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = h.store.Flight.UpdateFlight(ctx.Context(), filter, updateFlightParams)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Flight updated: " + flightID)
}

func (h *FlightHandler) HandleDeleteAllFlightsv1(ctx *fiber.Ctx) error {
	err := h.store.Flight.Drop(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Flights deleted")
}

func (h *FlightHandler) HandleGetSeatv1(ctx *fiber.Ctx) error {
	flightID := ctx.Params("fid")
	fid, err := primitive.ObjectIDFromHex(flightID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	seatID := ctx.Params("sid")
	sid, err := primitive.ObjectIDFromHex(seatID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	filter := bson.M{"_id": sid, "flight_id": fid}
	seat, err := h.store.Seat.GetSeat(ctx.Context(), filter)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(seat)
}
