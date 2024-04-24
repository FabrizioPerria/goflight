package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/types"
	"github.com/govalues/money"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SeedUsers(client *mongo.Client, dbName string) {
	userDb := db.NewMongoDbUserStore(client, dbName)
	userDb.Drop(context.Background())

	fmt.Println("Seeding users")
	for i := 0; i < 10; i++ {
		userParams := types.CreateUserParams{
			Email:         gofakeit.Email(),
			PlainPassword: gofakeit.Password(true, true, true, true, false, 10),
			Phone:         gofakeit.Phone(),
			FirstName:     gofakeit.FirstName(),
			LastName:      gofakeit.LastName(),
		}
		user, err := types.NewUserFromParams(userParams)
		if err != nil {
			fmt.Println(err)
			continue
		}
		userDb.CreateUser(context.Background(), user)
	}
}

var airports = []types.Airport{}

func SeedAirports(client *mongo.Client, dbName string) {
	airportDb := db.NewMongoDbAirportStore(client, dbName)
	airportDb.Drop(context.Background())

	fmt.Println("Seeding airports")
	for i := 0; i < 10; i++ {
		airport := types.Airport{
			City: gofakeit.City(),
			Code: gofakeit.DigitN(4),
		}
		airportDb.CreateAirport(context.Background(), airport)
		airports = append(airports, airport)
	}
}

func SeedFlights(client *mongo.Client, dbName string) {
	flightDb := db.NewMongoDbFlightStore(client, dbName)
	flightDb.Drop(context.Background())

	fmt.Println("Seeding flights")
	for i := 0; i < 10; i++ {
		departureKey := rand.Intn(len(airports))
		arrivalKey := rand.Intn(len(airports))

		flightParams := types.CreateFlightParams{
			Airline:   gofakeit.Company(),
			Departure: airports[departureKey],
			Arrival:   airports[arrivalKey],
			DepartureTime: types.FlightTime{
				Day:   gofakeit.Day(),
				Month: gofakeit.Month(),
				Year:  gofakeit.Year(),
				Hour:  gofakeit.Hour(),
				Min:   gofakeit.Minute(),
			},
			ArrivalTime: types.FlightTime{
				Day:   gofakeit.Day(),
				Month: gofakeit.Month(),
				Year:  gofakeit.Year(),
				Hour:  gofakeit.Hour(),
				Min:   gofakeit.Minute(),
			},
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
		fid, err := primitive.ObjectIDFromHex(newflight.Id)
		seats := []types.Seat{}
		for j := 0; j < 100; j++ {
			price, _ := money.NewAmountFromFloat64("USD", gofakeit.Float64Range(10, 1000))
			price = price.Round(2)
			priceFloat, _ := price.Float64()

			seat := types.Seat{
				Price:     priceFloat,
				Number:    j,
				Class:     types.SeatClass(gofakeit.Number(1, 3)),
				Available: true,
				FlightId:  fid.String(),
			}
			seats = append(seats, seat)
		}
		updateFlightParams := types.UpdateFlightParams{
			Seats: seats,
		}

		filter := bson.M{"_id": fid}
		_, err = flightDb.UpdateFlight(context.Background(), filter, updateFlightParams)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

const (
	uri    = "mongodb://localhost:27017"
	dbName = "goflight"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	SeedUsers(client, dbName)
	SeedAirports(client, dbName)
	SeedFlights(client, dbName)
}
