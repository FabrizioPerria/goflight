package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fabrizioperria/goflight/db"
	"github.com/fabrizioperria/goflight/types"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testFlightDb struct {
	Store  db.Store
	Client *mongo.Client
}

func setupFlightDb() (*testFlightDb, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	flightStore := db.NewMongoDbFlightStore(client)
	seatStore := db.NewMongoDbSeatStore(client, *flightStore)
	store := db.Store{Flight: flightStore, Seat: seatStore}
	return &testFlightDb{
		Client: client,
		Store:  store,
	}, nil
}

func teardownFlightDb(t *testing.T, testDb *testFlightDb) {
	if err := testDb.Store.Flight.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := testDb.Store.Seat.Drop(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := testDb.Client.Disconnect(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func getValidFlight() types.CreateFlightParams {
	return types.CreateFlightParams{
		Airline:       "Delta",
		Departure:     "JFK",
		Arrival:       "LAX",
		DepartureTime: "2021-01-01T00:00:00Z",
		ArrivalTime:   "2021-01-01T08:00:00Z",
	}
}

func createflight(flightHandler *FlightHandler, app *fiber.App, flight types.CreateFlightParams) (*http.Response, error) {
	app.Post("/api/v1/flights", flightHandler.HandlePostCreateFlightv1)
	flightMarshal, err := json.Marshal(flight)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest("POST", "/api/v1/flights", bytes.NewReader(flightMarshal))
	req.Header.Add("Content-Type", "application/json")
	return app.Test(req)
}

func getFlights(flightHandler *FlightHandler, app *fiber.App) (*http.Response, error) {
	app.Get("/api/v1/flights", flightHandler.HandleGetFlightsv1)

	req := httptest.NewRequest("GET", "/api/v1/flights", nil)
	req.Header.Add("Content-Type", "application/json")
	return app.Test(req)
}

func TestPostCreateFlightv1(t *testing.T) {
	db, err := setupFlightDb()
	assert.NoError(t, err)
	defer teardownFlightDb(t, db)
	flightHandler := FlightHandler{store: db.Store}

	app := fiber.New()

	flight := getValidFlight()
	response, error := createflight(&flightHandler, app, flight)
	assert.NoError(t, error)
	assert.Equal(t, 201, response.StatusCode)

	body, err := io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT := types.Flight{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	assert.NotEmpty(t, bodyT.Id)
	assert.Equal(t, flight.Airline, bodyT.Airline)
	assert.Equal(t, flight.Departure, bodyT.Departure)
	assert.Equal(t, flight.Arrival, bodyT.Arrival)
	assert.Equal(t, flight.DepartureTime, bodyT.DepartureTime)
	assert.Equal(t, flight.ArrivalTime, bodyT.ArrivalTime)

	// checke seats
	assert.Len(t, bodyT.Seats, 50)
	for i, seatId := range bodyT.Seats {
		filter := bson.M{"_id": seatId}
		seat, _ := db.Store.Seat.GetSeatById(context.Background(), filter)
		assert.Equal(t, bodyT.Id, seat.FlightId)
		assert.Equal(t, i, seat.Number)
		assert.True(t, seat.Available)
		assert.NotEqual(t, 0, seat.Price)
		assert.LessOrEqual(t, 1, seat.Class)
		assert.LessOrEqual(t, 1, seat.Location)
	}
}

func TestGetFlightsv1(t *testing.T) {
	db, err := setupFlightDb()
	assert.NoError(t, err)
	defer teardownFlightDb(t, db)
	flightHandler := FlightHandler{store: db.Store}

	app := fiber.New()

	flight := getValidFlight()
	_, error := createflight(&flightHandler, app, flight)
	assert.NoError(t, error)

	response, error := getFlights(&flightHandler, app)
	assert.NoError(t, error)
	assert.Equal(t, 200, response.StatusCode)

	body, err := io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT := []types.Flight{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	assert.Len(t, bodyT, 1)
}

func TestDeleteAllFlightsv1(t *testing.T) {
	db, err := setupFlightDb()
	assert.NoError(t, err)
	defer teardownFlightDb(t, db)
	flightHandler := FlightHandler{store: db.Store}

	app := fiber.New()

	flight := getValidFlight()
	_, error := createflight(&flightHandler, app, flight)
	assert.NoError(t, error)

	app.Delete("/api/v1/flights", flightHandler.HandleDeleteAllFlightsv1)

	req := httptest.NewRequest("DELETE", "/api/v1/flights", nil)
	req.Header.Add("Content-Type", "application/json")
	response, error := app.Test(req)
	assert.NoError(t, error)
	assert.Equal(t, 200, response.StatusCode)

	response, error = getFlights(&flightHandler, app)
	assert.NoError(t, error)
	assert.Equal(t, 200, response.StatusCode)

	body, err := io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT := []types.Flight{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	assert.Len(t, bodyT, 0)
}

func TestGetFlightByIdv1(t *testing.T) {
	db, err := setupFlightDb()
	assert.NoError(t, err)
	defer teardownFlightDb(t, db)
	flightHandler := FlightHandler{store: db.Store}

	app := fiber.New()

	flight := getValidFlight()
	response, error := createflight(&flightHandler, app, flight)
	assert.NoError(t, error)
	body, err := io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT := types.Flight{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)

	app.Get("/api/v1/flights/:fid", flightHandler.HandleGetFlightByIdv1)

	req := httptest.NewRequest("GET", "/api/v1/flights/"+bodyT.Id.Hex(), nil)
	req.Header.Add("Content-Type", "application/json")
	response, error = app.Test(req)
	assert.NoError(t, error)
	assert.Equal(t, 200, response.StatusCode)

	body, err = io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT = types.Flight{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	assert.NotEmpty(t, bodyT.Id)
	assert.Equal(t, flight.Airline, bodyT.Airline)
	assert.Equal(t, flight.Departure, bodyT.Departure)
	assert.Equal(t, flight.Arrival, bodyT.Arrival)
	assert.Equal(t, flight.DepartureTime, bodyT.DepartureTime)
	assert.Equal(t, flight.ArrivalTime, bodyT.ArrivalTime)
}

func TestPutFlightv1(t *testing.T) {
	db, err := setupFlightDb()
	assert.NoError(t, err)
	defer teardownFlightDb(t, db)
	flightHandler := FlightHandler{store: db.Store}

	app := fiber.New()

	flight := getValidFlight()
	response, error := createflight(&flightHandler, app, flight)
	assert.NoError(t, error)
	body, err := io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT := types.Flight{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	id := bodyT.Id.Hex()

	updateFlight := types.UpdateFlightParams{
		DepartureTime: "2024-01-01T01:00:00Z",
	}
	flight.DepartureTime = updateFlight.DepartureTime

	flightMarshal, err := json.Marshal(updateFlight)
	assert.NoError(t, err)

	req := httptest.NewRequest("PUT", "/api/v1/flights/"+id, bytes.NewReader(flightMarshal))
	req.Header.Add("Content-Type", "application/json")
	app.Put("/api/v1/flights/:fid", flightHandler.HandlePutFlightv1)

	response, error = app.Test(req)
	assert.NoError(t, error)
	assert.Equal(t, 200, response.StatusCode)

	req = httptest.NewRequest("GET", "/api/v1/flights/"+id, nil)
	req.Header.Add("Content-Type", "application/json")
	app.Get("/api/v1/flights/:fid", flightHandler.HandleGetFlightByIdv1)
	response, error = app.Test(req)
	assert.NoError(t, error)
	assert.Equal(t, 200, response.StatusCode)

	body, err = io.ReadAll(io.Reader(response.Body))
	assert.NoError(t, err)
	bodyT = types.Flight{}
	err = json.Unmarshal(body, &bodyT)
	assert.NoError(t, err)
	assert.NotEmpty(t, bodyT.Id)
	assert.Equal(t, flight.Airline, bodyT.Airline)
	assert.Equal(t, flight.Departure, bodyT.Departure)
	assert.Equal(t, flight.Arrival, bodyT.Arrival)
	assert.Equal(t, flight.DepartureTime, bodyT.DepartureTime)
	assert.Equal(t, flight.ArrivalTime, bodyT.ArrivalTime)
}
