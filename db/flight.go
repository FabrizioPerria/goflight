package db

import (
	"context"
	"fmt"

	"github.com/fabrizioperria/goflight/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FlightStorer interface {
	CreateFlight(ctx context.Context, flight *types.Flight) (*types.Flight, error)
	GetFlightById(ctx context.Context, id string) (*types.Flight, error)
	GetFlights(ctx context.Context) ([]*types.Flight, error)
	UpdateFlight(ctx context.Context, filter bson.M, values types.UpdateFlightParams) (string, error)
	Dropper
}

const (
	flightCollection = "flights"
)

type MongoDbFlightStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDbFlightStore(client *mongo.Client) *MongoDbFlightStore {
	return &MongoDbFlightStore{
		client:     client,
		collection: client.Database(DBNAME).Collection(flightCollection),
	}
}

func (db *MongoDbFlightStore) GetFlightById(ctx context.Context, id string) (*types.Flight, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var flight types.Flight
	err = db.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&flight)
	if err != nil {
		return nil, err
	}
	return &flight, nil
}

func (db *MongoDbFlightStore) GetFlights(ctx context.Context) ([]*types.Flight, error) {
	var cursor *mongo.Cursor
	cursor, err := db.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	results := make([]*types.Flight, 0)
	err = cursor.All(ctx, &results)

	return results, err
}

func (db *MongoDbFlightStore) CreateFlight(ctx context.Context, flight *types.Flight) (*types.Flight, error) {
	result, err := db.collection.InsertOne(ctx, flight)
	flight.Id = result.InsertedID.(primitive.ObjectID)
	return flight, err
}

func (db *MongoDbFlightStore) Drop(ctx context.Context) error {
	return db.collection.Drop(ctx)
}

func (db *MongoDbFlightStore) UpdateFlight(ctx context.Context, filter bson.M, values types.UpdateFlightParams) (string, error) {
	thinValues := bson.M{}
	if values.ArrivalTime != "" {
		thinValues["arrival_time"] = values.ArrivalTime
	}
	if values.DepartureTime != "" {
		thinValues["departure_time"] = values.DepartureTime
	}
	if len(values.Seats) > 0 {
		thinValues["seats"] = values.Seats
	}
	update := bson.M{"$set": thinValues}
	result, err := db.collection.UpdateOne(ctx, filter, update)
	if err != nil || result.ModifiedCount == 0 {
		return "", fmt.Errorf("flight not found")
	}
	return "", nil
}
