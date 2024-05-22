package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reservation struct {
	ReservationDate string             `json:"reservation_date,omitempty" bson:"reservation_date,omitempty"`
	Id              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SeatId          primitive.ObjectID `json:"seat_id,omitempty" bson:"seat_id,omitempty"`
	UserId          primitive.ObjectID `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

type CreateReservationParams struct {
	SeatId primitive.ObjectID `json:"seat_id" bson:"seat_id"`
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
}

func ReservationFromParams(params *CreateReservationParams) *Reservation {
	return &Reservation{
		SeatId:          params.SeatId,
		UserId:          params.UserId,
		ReservationDate: time.Now().Format(time.RFC3339),
	}
}
