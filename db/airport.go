package db

import (
	"context"

	"github.com/fabrizioperria/goflight/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type AirportStorer interface {
	Dropper
}

type MongoDbAirportStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

const airportCollection = "airports"

func NewMongoDbAirportStore(client *mongo.Client, dbName string) *MongoDbAirportStore {
	return &MongoDbAirportStore{
		client:     client,
		collection: client.Database(dbName).Collection(airportCollection),
	}
}

func (db *MongoDbAirportStore) Drop(ctx context.Context) error {
	return db.collection.Drop(ctx)
}

func (db *MongoDbAirportStore) CreateAirport(ctx context.Context, airport types.Airport) error {
	_, err := db.collection.InsertOne(ctx, airport)
	return err
}
