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
	CreateFlight(ctx context.Context, user *types.User) (*types.User, error)
	GetFlightById(ctx context.Context, id string) (*types.User, error)
	GetFlights(ctx context.Context) ([]*types.User, error)
	DeleteFlightById(ctx context.Context, id string) (string, error)
	UpdateFlight(ctx context.Context, filter bson.M, values types.UpdateUserParams) (string, error)
	Dropper
}

const (
	flightCollection = "flights"
)

type MongoDbFlightStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDbFlightStore(client *mongo.Client, dbName string) *MongoDbFlightStore {
	return &MongoDbFlightStore{
		client:     client,
		collection: client.Database(dbName).Collection(flightCollection),
	}
}

func (db *MongoDbFlightStore) CreateFlight(ctx context.Context, flight *types.Flight) (*types.Flight, error) {
	result, err := db.collection.InsertOne(ctx, flight)
	flight.Id = result.InsertedID.(primitive.ObjectID).Hex()
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
