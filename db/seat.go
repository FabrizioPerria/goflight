package db

import (
	"context"
	"fmt"

	"github.com/fabrizioperria/goflight/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type SeatStorer interface {
	CreateSeat(ctx context.Context, user *types.Seat) (*types.Seat, error)
	UpdateSeat(ctx context.Context, filter bson.M, values types.UpdateSeatParams) (string, error)
	GetSeats(ctx context.Context, filter bson.M) ([]*types.Seat, error)
	GetSeat(ctx context.Context, filter bson.M) (*types.Seat, error)
	ReserveSeat(ctx context.Context, filter bson.M, userId primitive.ObjectID) (*types.Reservation, error)
	Dropper
}

const (
	seatCollection = "seats"
)

type MongoDbSeatStore struct {
	client           *mongo.Client
	collection       *mongo.Collection
	flightStore      MongoDbFlightStore
	reservationStore MongoDbReservationStore
}

func NewMongoDbSeatStore(client *mongo.Client, flightStore MongoDbFlightStore, reservationStore MongoDbReservationStore) *MongoDbSeatStore {
	return &MongoDbSeatStore{
		client:           client,
		collection:       client.Database(DBNAME).Collection(seatCollection),
		flightStore:      flightStore,
		reservationStore: reservationStore,
	}
}

func (db *MongoDbSeatStore) CreateSeat(ctx context.Context, seat *types.Seat) (*types.Seat, error) {
	result, err := db.collection.InsertOne(ctx, seat)
	if err != nil {
		return nil, err
	}
	seat.Id = result.InsertedID.(primitive.ObjectID)

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

	return "", nil
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

func (db *MongoDbSeatStore) GetSeat(ctx context.Context, filter bson.M) (*types.Seat, error) {
	var seat types.Seat
	err := db.collection.FindOne(ctx, filter).Decode(&seat)
	if err != nil {
		return nil, err
	}
	return &seat, nil
}

func (db *MongoDbSeatStore) ReserveSeat(ctx context.Context, filter bson.M, userId primitive.ObjectID) (*types.Reservation, error) {
	session, err := db.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		seat, err := db.GetSeat(sessionContext, filter)
		if err != nil {
			return nil, err
		}

		if !seat.Available {
			return nil, fmt.Errorf("seat not available")
		}

		seat.Available = false
		_, err = db.UpdateSeat(sessionContext, filter, types.UpdateSeatParams{
			Available: seat.Available,
			Price:     seat.Price,
		})
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("fake error")

		filter = bson.M{"_id": seat.FlightId}
		update := bson.M{"$pull": bson.M{"seats": seat.Id}}
		_, err = db.flightStore.collection.UpdateOne(sessionContext, filter, update)
		if err != nil {
			return nil, err
		}

		reservationParams := types.CreateReservationParams{
			UserId: userId,
			SeatId: seat.Id,
		}
		reservation := types.ReservationFromParams(&reservationParams)
		result, err := db.reservationStore.CreateReservation(sessionContext, reservation)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	reservation, err := session.WithTransaction(ctx, callback, txnOpts)
	if err != nil {
		return nil, err
	}

	return reservation.(*types.Reservation), nil
}
