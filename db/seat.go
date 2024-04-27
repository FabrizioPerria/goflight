package db

import (
	"context"

	"github.com/fabrizioperria/goflight/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SeatStorer interface {
	CreateSeat(ctx context.Context, user *types.Seat) (*types.Seat, error)
	UpdateSeat(ctx context.Context, filter bson.M, values types.UpdateSeatParams) (string, error)
	GetSeats(ctx context.Context, filter bson.M) ([]*types.Seat, error)
	Dropper
}

const (
	seatCollection = "seats"
)

type MongoDbSeatStore struct {
	client      *mongo.Client
	collection  *mongo.Collection
	flightStore MongoDbFlightStore
}

func NewMongoDbSeatStore(client *mongo.Client, flightStore MongoDbFlightStore) *MongoDbSeatStore {
	return &MongoDbSeatStore{
		client:      client,
		collection:  client.Database(DBNAME).Collection(seatCollection),
		flightStore: flightStore,
	}
}

func (db *MongoDbSeatStore) CreateSeat(ctx context.Context, seat *types.Seat) (*types.Seat, error) {
	result, err := db.collection.InsertOne(ctx, seat)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": result.InsertedID}
	update := bson.M{"$addToSet": bson.M{"seats": seat}}
	_, err = db.flightStore.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return seat, err
}

func (db *MongoDbSeatStore) UpdateSeat(ctx context.Context, filter bson.M, values types.UpdateSeatParams) (string, error) {
	update := bson.M{"$set": values}
	result, err := db.collection.UpdateOne(ctx, filter, update)
	if err != nil || result.ModifiedCount == 0 {
		return "", err
	}

	return result.UpsertedID.(primitive.ObjectID).Hex(), err
}

func (db *MongoDbSeatStore) Drop(ctx context.Context) error {
	return db.collection.Drop(ctx)
}

func (db *MongoDbSeatStore) GetSeats(ctx context.Context, filter bson.M) ([]*types.Seat, error) {
	var cursor *mongo.Cursor
	cursor, err := db.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	results := make([]*types.Seat, 0)
	err = cursor.All(ctx, &results)

	return results, err
}
