package db

import (
	"context"

	"github.com/fabrizioperria/goflight/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SeatStorer interface {
	CreateSeat(ctx context.Context, user *types.User) (*types.User, error)
	UpdateSeat(ctx context.Context, filter bson.M, values types.UpdateUserParams) (string, error)
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
	seat.Id = result.InsertedID.(primitive.ObjectID).Hex()

	fid, _ := primitive.ObjectIDFromHex(seat.FlightId)

	filter := bson.M{"_id": fid}
	update := bson.M{"$addToSet": bson.M{"seats": seat}}
	_, err = db.flightStore.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return seat, err
}

func (db *MongoDbSeatStore) UpdateSeat(ctx context.Context, filter bson.M, values types.UpdateSeatParams) (string, error) {
	// update := bson.M{"$set": values}
	// result, err := db.collection.UpdateOne(ctx, filter, update)
	// if err != nil || result.ModifiedCount == 0 {
	// 	return "", fmt.Errorf("seat not found")
	// }
	//
	// update = bson.M{"$set": bson.M{"seats.$[elem].status": values.Status}}
	//
	// result, err = db.flightCollection.UpdateOne(ctx, filter, update)
	// if err != nil || result.ModifiedCount == 0 {
	// 	return "", fmt.Errorf("flight not found")
	// }

	return "", nil
}

func (db *MongoDbSeatStore) Drop(ctx context.Context) error {
	return db.collection.Drop(ctx)
}
