package db

import (
	"context"

	"github.com/fabrizioperria/goflight/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ReservationStorer interface {
	CreateReservation(ctx context.Context, user *types.Reservation) (*types.Reservation, error)
	GetReservations(ctx context.Context) ([]*types.Reservation, error)
	GetReservation(ctx context.Context, filter bson.M) (*types.Reservation, error)
	DeleteReservation(ctx context.Context, filter bson.M) error
	Dropper
}

type MongoDbReservationStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

const (
	reservationCollection = "reservations"
)

func NewMongoDbReservationStore(client *mongo.Client) *MongoDbReservationStore {
	return &MongoDbReservationStore{
		client:     client,
		collection: client.Database(DBNAME).Collection(reservationCollection),
	}
}

func (db *MongoDbReservationStore) CreateReservation(ctx context.Context, reservation *types.Reservation) (*types.Reservation, error) {
	result, err := db.collection.InsertOne(ctx, reservation)
	if err != nil {
		return nil, err
	}
	reservation.Id = result.InsertedID.(primitive.ObjectID)

	return reservation, err
}

func (db *MongoDbReservationStore) GetReservations(ctx context.Context) ([]*types.Reservation, error) {
	var reservations []*types.Reservation
	cursor, err := db.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (db *MongoDbReservationStore) GetReservation(ctx context.Context, filter bson.M) (*types.Reservation, error) {
	var reservation *types.Reservation
	err := db.collection.FindOne(ctx, filter).Decode(&reservation)
	if err != nil {
		return nil, err
	}
	return reservation, nil
}

func (db *MongoDbReservationStore) DeleteReservation(ctx context.Context, filter bson.M) error {
	_, err := db.collection.DeleteOne(ctx, filter)
	return err
}

func (db *MongoDbReservationStore) Drop(ctx context.Context) error {
	return db.collection.Drop(ctx)
}
