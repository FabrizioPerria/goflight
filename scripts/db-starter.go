package main

import (
	"context"
	"fmt"
	"log"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SeedUsers(client *mongo.Client, dbName string) {
	userDb := db.NewMongoDbUserStore(client, dbName)
	userDb.Drop(context.Background())

	fmt.Println("Seeding users")
	for i := 0; i < 100; i++ {
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

var airports = map[string]types.Airport{}

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
		airports[airport.Code] = airport
	}
}

func SeedFlights(client *mongo.Client, dbName string) {
	// flightDb := db.NewMongoDbFlightStore(client, dbName)
	// flightDb.Drop(context.Background())
	//
	// keys := make([]string, 0, len(airports))
	// for k := range airports {
	// 	keys = append(keys, k)
	// }
	//
	// fmt.Println("Seeding flights")
	// for i := 0; i < 10; i++ {
	// 	departureKey := keys[gofakeit.Number(0, len(keys)-1)]
	// 	arrivalKey := keys[gofakeit.Number(0, len(keys)-1)]
	//
	// 	flightParams := types.CreateFlightParams{
	// 		Airline:   gofakeit.Company(),
	// 		Departure: airports[departureKey],
	// 		Arrival:   airports[arrivalKey],
	// 		DepartureTime: types.FlightTime{
	// 			Day:   gofakeit.Day(),
	// 			Month: gofakeit.Month(),
	// 			Year:  gofakeit.Year(),
	// 			Hour:  gofakeit.Hour(),
	// 			Min:   gofakeit.Minute(),
	// 		},
	// 		ArrivalTime: types.FlightTime{
	// 			Day:   gofakeit.Day(),
	// 			Month: gofakeit.Month(),
	// 			Year:  gofakeit.Year(),
	// 			Hour:  gofakeit.Hour(),
	// 			Min:   gofakeit.Minute(),
	// 		},
	// 	}
	// 	// flight, err := types.NewFlightFromParams(flightParams)
	// 	// if err != nil {
	// 	// 	fmt.Println(err)
	// 	// 	continue
	// 	// }
	// 	// flightDb.CreateFlight(context.Background(), flight)
	// }
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
