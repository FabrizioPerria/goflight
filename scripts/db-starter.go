package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/types"
	"github.com/govalues/money"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	users        = []types.User{}
	flights      = []types.Flight{}
	reservations = []types.Reservation{}
)

func SeedUsers(client *mongo.Client) {
	userDb := db.NewMongoDbUserStore(client)
	userDb.Drop(context.Background())

	fmt.Println("Seeding users")
	userParams := types.CreateUserParams{
		Email:         "a.b@c.d",
		PlainPassword: "password",
		Phone:         "1234567890",
		FirstName:     "Dude",
		LastName:      "Dudely",
	}
	user, err := types.NewUserFromParams(userParams, true)
	if err != nil {
		fmt.Println(err)
	}
	userDb.CreateUser(context.Background(), user)
	users = append(users, *user)
	userParams = types.CreateUserParams{
		Email:         "nonadmin@a.b",
		PlainPassword: "password",
		Phone:         "1234567890",
		FirstName:     "Dude",
		LastName:      "Dudely",
	}
	user, err = types.NewUserFromParams(userParams, false)
	if err != nil {
		fmt.Println(err)
	}
	userDb.CreateUser(context.Background(), user)
	users = append(users, *user)

	for i := 0; i < 10; i++ {
		userParams := types.CreateUserParams{
			Email:         gofakeit.Email(),
			PlainPassword: gofakeit.Password(true, true, true, true, false, 10),
			Phone:         gofakeit.Phone(),
			FirstName:     gofakeit.FirstName(),
			LastName:      gofakeit.LastName(),
		}
		user, err := types.NewUserFromParams(userParams, false)
		if err != nil {
			fmt.Println(err)
			continue
		}
		userDb.CreateUser(context.Background(), user)
		users = append(users, *user)
	}
}

func SeedFlights(client *mongo.Client) {
	flightDb := db.NewMongoDbFlightStore(client)
	flightDb.Drop(context.Background())
	reservationDb := db.NewMongoDbReservationStore(client)
	reservationDb.Drop(context.Background())
	seatDb := db.NewMongoDbSeatStore(client, *flightDb, *reservationDb)
	seatDb.Drop(context.Background())

	fmt.Println("Seeding flights")
	for i := 0; i < 10; i++ {

		flightParams := types.CreateFlightParams{
			Airline:       gofakeit.Company(),
			Departure:     gofakeit.City(),
			Arrival:       gofakeit.City(),
			DepartureTime: gofakeit.Date().Format(time.RFC3339),
			ArrivalTime:   gofakeit.Date().Format(time.RFC3339),
			NumberOfSeats: gofakeit.Number(10, 100),
		}
		flight, err := types.NewFlightFromParams(flightParams)
		if err != nil {
			fmt.Println(err)
			continue
		}
		newflight, err := flightDb.CreateFlight(context.Background(), flight)
		if err != nil {
			fmt.Println(err)
			continue
		}
		updateData := types.UpdateFlightParams{
			Seats:         []primitive.ObjectID{},
			ArrivalTime:   newflight.ArrivalTime,
			DepartureTime: newflight.DepartureTime,
		}
		for j := 0; j < flightParams.NumberOfSeats; j++ {
			price, _ := money.NewAmountFromFloat64("USD", gofakeit.Float64Range(10, 1000))
			price = price.Round(2)
			priceFloat, _ := price.Float64()

			seat := types.Seat{
				Price:     priceFloat,
				Number:    j,
				Class:     types.SeatClass(gofakeit.Number(1, 3)),
				Location:  types.SeatLocation(gofakeit.Number(1, 3)),
				Available: true,
				FlightId:  newflight.Id,
			}
			registeredSeat, err := seatDb.CreateSeat(context.Background(), &seat)
			if err != nil {
				fmt.Println(err)
				continue
			}
			updateData.Seats = append(updateData.Seats, registeredSeat.Id)
		}

		filter := bson.M{"_id": newflight.Id}
		_, err = flightDb.UpdateFlight(context.Background(), filter, updateData)
		if err != nil {
			fmt.Println(err)
			continue
		}
		newflight.Seats = updateData.Seats
		flights = append(flights, *newflight)
	}
}

func SeedReservations(client *mongo.Client) {
	reservationDb := db.NewMongoDbReservationStore(client)
	reservationDb.Drop(context.Background())

	fmt.Println("Seeding reservations")
	for i := 0; i < 10; i++ {
		flight := flights[gofakeit.Number(0, len(flights)-1)]
		seat := flight.Seats[gofakeit.Number(0, len(flight.Seats)-1)]
		user := users[gofakeit.Number(0, len(users)-1)]
		reservation := &types.Reservation{
			UserId: user.Id,
			SeatId: seat,
		}
		_, err := reservationDb.CreateReservation(context.Background(), reservation)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

const (
	uri = "mongodb://mongo1:27017"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetReplicaSet("rs0"))
	if err != nil {
		log.Fatal(err)
	}

	SeedUsers(client)
	SeedFlights(client)
	SeedReservations(client)
}
