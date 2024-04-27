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
	// GetFlightById(ctx context.Context, id string) (*types.Flight, error)
	GetFlights(ctx context.Context) ([]*types.Flight, error)
	// DeleteFlightById(ctx context.Context, id string) (string, error)
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
	update := bson.M{"$set": values}
	result, err := db.collection.UpdateOne(ctx, filter, update)
	if err != nil || result.ModifiedCount == 0 {
		return "", fmt.Errorf("flight not found")
	}
	return "", nil
}
