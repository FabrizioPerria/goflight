package db

import (
	"context"

	"github.com/fabrizioperria/goflight/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SeatStorer interface {
	CreateSeat(ctx context.Context, user *types.Seat) (*types.Seat, error)
	UpdateSeat(ctx context.Context, filter Map, values types.UpdateSeatParams) (string, error)
	GetSeats(ctx context.Context, filter Map, pagination *Pagination) ([]*types.Seat, error)
	GetSeat(ctx context.Context, filter Map) (*types.Seat, error)
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
	dbName := os.Getenv("DB_NAME")
	return &MongoDbSeatStore{
		client:      client,
		collection:  client.Database(dbName).Collection(seatCollection),
		flightStore: flightStore,
	}
}

func (db *MongoDbSeatStore) CreateSeat(ctx context.Context, seat *types.Seat) (*types.Seat, error) {
	result, err := db.collection.InsertOne(ctx, seat)
	if err != nil {
		return nil, err
	}
	seat.Id = result.InsertedID.(primitive.ObjectID)

	filter := Map{"_id": result.InsertedID}
	update := Map{"$addToSet": Map{"seats": seat}}
	_, err = db.flightStore.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return seat, err
}

func (db *MongoDbSeatStore) UpdateSeat(ctx context.Context, filter Map, values types.UpdateSeatParams) (string, error) {
	update := Map{"$set": values}
	result, err := db.collection.UpdateOne(ctx, filter, update)
	if err != nil || result.ModifiedCount == 0 {
		return "", err
	}

	return "", nil
}

func (db *MongoDbSeatStore) Drop(ctx context.Context) error {
	return db.collection.Drop(ctx)
}

func (db *MongoDbSeatStore) GetSeats(ctx context.Context, filter Map, pagination *Pagination) ([]*types.Seat, error) {
	var cursor *mongo.Cursor
	cursor, err := db.collection.Find(ctx, filter, pagination.ToFindOptions())
	if err != nil {
		return nil, err
	}

	results := make([]*types.Seat, 0)
	err = cursor.All(ctx, &results)

	return results, err
}

func (db *MongoDbSeatStore) GetSeat(ctx context.Context, filter Map) (*types.Seat, error) {
	var seat types.Seat
	err := db.collection.FindOne(ctx, filter).Decode(&seat)
	if err != nil {
		return nil, err
	}
	return &seat, nil
}
