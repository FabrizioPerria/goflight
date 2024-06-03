package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/db/fixtures"
	"github.com/fabrizioperria/goflight/types"
	"github.com/govalues/money"
	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	users        = []types.User{}
	flights      = []types.Flight{}
	reservations = []types.Reservation{}
)

func SeedUsers(client *mongo.Client, store *db.Store) {
	store.User.Drop(context.Background())

	fmt.Println("Seeding users")
	user, err := fixtures.AddUser(store, "admin@a.b", "password", "1234567890", "Dude", "Dudely", true)
	if err != nil {
		fmt.Println(err)
	}
	users = append(users, *user)
	user, _ = fixtures.AddUser(store, "nonadmin@a.b", "password", "1234567890", "Dude", "Dudely", false)
	users = append(users, *user)

	for i := 0; i < 10; i++ {
		user, err := fixtures.AddUser(store,
			gofakeit.Email(),
			gofakeit.Password(true, true, true, true, false, 10),
			gofakeit.Phone(),
			gofakeit.FirstName(),
			gofakeit.LastName(),
			false)
		if err != nil {
			fmt.Println(err)
			continue
		}

		users = append(users, *user)
	}
}

func SeedFlights(client *mongo.Client, store *db.Store) {
	store.Flight.Drop(context.Background())
	store.Seat.Drop(context.Background())
	store.Reservation.Drop(context.Background())

	fmt.Println("Seeding flights")
	for i := 0; i < 10; i++ {

		numberSeats := gofakeit.Number(1, 100)
		newflight, _ := fixtures.AddFlight(store,
			gofakeit.Company(),
			gofakeit.City(),
			gofakeit.City(),
			time.Now().AddDate(0, 0, gofakeit.Number(1, 30)).Format(time.RFC3339),
			time.Now().AddDate(0, 0, gofakeit.Number(31, 60)).Format(time.RFC3339),
			numberSeats)

		newSeats := []primitive.ObjectID{}
		for j := 0; j < numberSeats; j++ {
			price, _ := money.NewAmountFromFloat64("USD", gofakeit.Float64Range(10, 1000))
			seat, _ := fixtures.AddSeat(store,
				price,
				j,
				types.SeatClass(gofakeit.Number(0, 2)),
				types.SeatLocation(gofakeit.Number(0, 2)),
				true,
				newflight.Id)
			newSeats = append(newSeats, seat.Id)
		}
		err := fixtures.AddSeatsToFlight(store, newflight.Id, newSeats)
		if err != nil {
			fmt.Println(err)
			// should delete the flight, but oh well
			continue
		}
		newflight.Seats = newSeats

		flights = append(flights, *newflight)
	}
}

func SeedReservations(client *mongo.Client, store *db.Store) {
	store.Reservation.Drop(context.Background())

	fixtures.AddReservation(store, flights[0].Seats[0], users[1].Id)

	fmt.Println("Seeding reservations")
	for i := 0; i < 10; i++ {
		flight := flights[gofakeit.Number(0, len(flights)-1)]
		seat := flight.Seats[gofakeit.Number(0, len(flight.Seats)-1)]
		user := users[gofakeit.Number(0, len(users)-1)]
		fixtures.AddReservation(store, seat, user.Id)
	}
}

func main() {
	mongoUrl := os.Getenv("MONGO_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUrl).SetReplicaSet("rs0"))
	if err != nil {
		log.Fatal(err)
	}

	userDb := db.NewMongoDbUserStore(client)
	flightDb := db.NewMongoDbFlightStore(client)
	seatDb := db.NewMongoDbSeatStore(client, *flightDb)
	reservationDb := db.NewMongoDbReservationStore(client, *flightDb, *seatDb)

	store := db.NewStore(userDb, flightDb, seatDb, reservationDb)
	SeedUsers(client, store)
	SeedFlights(client, store)
	SeedReservations(client, store)
}
