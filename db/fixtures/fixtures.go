package fixtures

import (
	"context"

	"github.com/fabrizioperria/goflight/handlers/middleware"
	"github.com/fabrizioperria/goflight/types"

	"github.com/govalues/money"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/fabrizioperria/goflight/db"
)

func getValidUser() types.CreateUserParams {
	return types.CreateUserParams{
		FirstName:     "Frank",
		LastName:      "Potato",
		Email:         "fp@test.com",
		Phone:         "123456789",
		PlainPassword: "password",
	}
}

func AuthenticateUser(store *db.Store) (*types.User, string) {
	userParams := getValidUser()
	user, _ := AddUser(store, userParams.Email, userParams.PlainPassword, userParams.Phone, userParams.FirstName, userParams.LastName, true)

	token := middleware.ProduceToken(user)
	return user, token
}

func AddUser(store *db.Store, email, password, phone, firstName, lastName string, isAdmin bool) (*types.User, error) {
	userParams := types.CreateUserParams{
		Email:         email,
		PlainPassword: password,
		Phone:         phone,
		FirstName:     firstName,
		LastName:      lastName,
	}
	user, err := types.NewUserFromParams(userParams, isAdmin)
	if err != nil {
		return nil, err
	}
	return store.User.CreateUser(context.Background(), user)
}

func AddFlight(store *db.Store, airline, departure, arrival, departureTime, arrivalTime string, numSeats int) (*types.Flight, error) {
	flightParams := types.CreateFlightParams{
		Airline:       airline,
		Departure:     departure,
		Arrival:       arrival,
		DepartureTime: departureTime,
		ArrivalTime:   arrivalTime,
		NumberOfSeats: numSeats,
	}
	flight, err := types.NewFlightFromParams(flightParams)
	if err != nil {
		return nil, err
	}
	return store.Flight.CreateFlight(context.Background(), flight)
}

func AddSeat(store *db.Store, price money.Amount, number int, class types.SeatClass, location types.SeatLocation, available bool, flightId primitive.ObjectID) (*types.Seat, error) {
	price = price.Round(2)
	priceFloat, _ := price.Float64()

	seat := types.Seat{
		Price:     priceFloat,
		Number:    number,
		Class:     class,
		Location:  location,
		Available: available,
		FlightId:  flightId,
	}
	return store.Seat.CreateSeat(context.Background(), &seat)
}

func AddSeatsToFlight(store *db.Store, flightId primitive.ObjectID, seats []primitive.ObjectID) error {
	filter := db.Map{"_id": flightId}
	flight, err := store.Flight.GetFlight(context.Background(), filter)
	if err != nil {
		return err
	}
	updateData := types.UpdateFlightParams{
		Seats:         seats,
		ArrivalTime:   flight.ArrivalTime,
		DepartureTime: flight.DepartureTime,
	}

	_, err = store.Flight.UpdateFlight(context.Background(), filter, updateData)
	return err
}

func AddReservation(store *db.Store, seatId primitive.ObjectID, userId primitive.ObjectID) (*types.Reservation, error) {
	return store.Reservation.CreateReservation(context.Background(), db.Map{"_id": seatId}, userId)
}
