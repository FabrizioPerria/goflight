package db

import (
	"context"
	"fmt"
	"time"

	"github.com/fabrizioperria/goflight/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type ReservationStorer interface {
	CreateReservation(ctx context.Context, filter Map, userId primitive.ObjectID) (*types.Reservation, error)
	GetReservations(ctx context.Context, filter Map, pagination *Pagination) ([]*types.Reservation, error)
	GetReservation(ctx context.Context, filter Map) (*types.Reservation, error)
	DeleteReservation(ctx context.Context, filter Map) error
	Dropper
}

type MongoDbReservationStore struct {
	client      *mongo.Client
	collection  *mongo.Collection
	flightStore MongoDbFlightStore
	seatStore   MongoDbSeatStore
}

const (
	reservationCollection = "reservations"
)

func NewMongoDbReservationStore(client *mongo.Client, flightStore MongoDbFlightStore, seatStore MongoDbSeatStore) *MongoDbReservationStore {
	return &MongoDbReservationStore{
		client:      client,
		collection:  client.Database(DBNAME).Collection(reservationCollection),
		flightStore: flightStore,
		seatStore:   seatStore,
	}
}

func (db *MongoDbReservationStore) CreateReservation(ctx context.Context, filter Map, userId primitive.ObjectID) (*types.Reservation, error) {
	session, err := db.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		seat, err := db.seatStore.GetSeat(sessionContext, filter)
		if err != nil {
			return nil, err
		}

		if !seat.Available {
			return nil, fmt.Errorf("seat not available")
		}

		if _, err = db.seatStore.UpdateSeat(sessionContext, filter, types.UpdateSeatParams{Available: false, Price: seat.Price}); err != nil {
			return nil, err
		}

		filter = Map{"_id": seat.FlightId}
		update := Map{"$pull": Map{"seats": seat.Id}}
		_, err = db.flightStore.collection.UpdateOne(sessionContext, filter, update)
		if err != nil {
			return nil, err
		}

		reservationParams := types.CreateReservationParams{
			UserId: userId,
			SeatId: seat.Id,
		}
		reservation := types.ReservationFromParams(&reservationParams)

		reservation.ReservationDate = time.Now().Format(time.RFC3339)
		reservation.CancellationDate = ""
		result, err := db.collection.InsertOne(ctx, reservation)
		if err != nil {
			return nil, err
		}
		reservation.Id = result.InsertedID.(primitive.ObjectID)

		return reservation.Id, nil
	}

	reservationId, err := session.WithTransaction(ctx, callback, txnOpts)
	if err != nil {
		return nil, err
	}
	reservation, err := db.GetReservation(ctx, Map{"_id": reservationId.(primitive.ObjectID)})
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (db *MongoDbReservationStore) GetReservations(ctx context.Context, filter Map, pagination *Pagination) ([]*types.Reservation, error) {
	var reservations []*types.Reservation
	cursor, err := db.collection.Find(ctx, filter, pagination.ToFindOptions())
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &reservations); err != nil {
		return nil, err
	}
	return reservations, nil
}

func (db *MongoDbReservationStore) GetReservation(ctx context.Context, filter Map) (*types.Reservation, error) {
	var reservation *types.Reservation
	err := db.collection.FindOne(ctx, filter).Decode(&reservation)
	if err != nil {
		return nil, err
	}
	return reservation, nil
}

func (db *MongoDbReservationStore) DeleteReservation(ctx context.Context, filter Map) error {
	session, err := db.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)

	callback := func(sessionContext mongo.SessionContext) (interface{}, error) {
		reservation, err := db.GetReservation(ctx, filter)
		if err != nil {
			return nil, err
		}

		if reservation.CancellationDate != "" {
			return nil, fmt.Errorf("reservation already cancelled")
		}

		seatFilter := Map{"_id": reservation.SeatId}
		seat, err := db.seatStore.GetSeat(ctx, seatFilter)
		if err != nil {
			return nil, err
		}

		if _, err = db.seatStore.UpdateSeat(ctx, seatFilter, types.UpdateSeatParams{Available: true, Price: seat.Price}); err != nil {
			return nil, err
		}

		flightFilter := Map{"_id": seat.FlightId}
		update := Map{"$push": Map{"seats": seat.Id}}
		_, err = db.flightStore.collection.UpdateOne(sessionContext, flightFilter, update)
		if err != nil {
			return nil, err
		}

		result, err := db.collection.UpdateOne(ctx, filter, Map{"$set": Map{"cancellation_date": time.Now().Format(time.RFC3339)}})
		if err != nil {
			return nil, err
		}
		return result.UpsertedID, nil
	}
	_, err = session.WithTransaction(ctx, callback, txnOpts)
	return err
}

func (db *MongoDbReservationStore) Drop(ctx context.Context) error {
	return db.collection.Drop(ctx)
}
