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

func (h *FlightHandler) HandleGetFlightByIdv1(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	flight, err := h.store.Flight.GetFlightById(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(flight)
}

func (h *FlightHandler) HandleGetSeatsv1(ctx *fiber.Ctx) error {
	flightID := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(flightID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	filter := bson.M{"flight_id": oid}
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

	for i := 0; i < 50; i++ {
		seat := types.Seat{
			FlightId:  primitive.ObjectID(flight.Id),
			Number:    i,
			Price:     100,
			Class:     types.SeatClass(i%3 + 1),
			Available: true,
		}
		_, err = h.store.Seat.CreateSeat(ctx.Context(), &seat)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return ctx.Status(fiber.StatusCreated).JSON(flight)
}

func (h *FlightHandler) HandlePutFlightv1(ctx *fiber.Ctx) error {
	flightID := ctx.Params("id")
	updateFlightParams := types.UpdateFlightParams{}
	err := ctx.BodyParser(&updateFlightParams)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	filter := map[string]interface{}{"_id": flightID}
	_, err = h.store.Flight.UpdateFlight(ctx.Context(), filter, updateFlightParams)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(fiber.Map{"id": flightID})
}

func (h *FlightHandler) HandleDeleteAllFlightsv1(ctx *fiber.Ctx) error {
	err := h.store.Flight.Drop(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).SendString("Flights deleted")
}
