package handlers

import (
	"context"
	"testing"

	"github.com/fabrizioperria/goflight/db"
	_ "github.com/gofiber/fiber/v2"
	_ "github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testFlightDb struct {
	FlightStore db.FlightStorer
	SeatStore   db.SeatStorer
	Client      *mongo.Client
}

func setupFlightDb() (*testFlightDb, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	flightStore := db.NewMongoDbFlightStore(client)
	seatStore := db.NewMongoDbSeatStore(client, *flightStore)
	return &testFlightDb{
		Client:      client,
		FlightStore: flightStore,
		SeatStore:   seatStore,
	}, nil
}

func teardownFlightDb(t *testing.T, testDb *testFlightDb) {
	if err := testDb.FlightStore.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := testDb.Client.Disconnect(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestPostCreateFlightv1(t *testing.T) {
	// db, err := setupFlightDb()
	// assert.NoError(t, err)
	// defer teardownFlightDb(t, db)
	// flightHandler := FlightHandler{
	// 	FlightStore: db.FlightStore,
	// }
	//
	// app := fiber.New()
}
