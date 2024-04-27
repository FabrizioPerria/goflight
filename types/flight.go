package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Flight struct {
	Departure     string               `json:"departure" bson:"departure"`
	Arrival       string               `json:"arrival" bson:"arrival"`
	Airline       string               `json:"airline" bson:"airline"`
	DepartureTime string               `json:"departure_time" bson:"departure_time"`
	ArrivalTime   string               `json:"arrival_time" bson:"arrival_time"`
	Seats         []primitive.ObjectID `json:"seats" bson:"seats"`
	Id            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
}

type CreateFlightParams struct {
	Departure     string `json:"departure" bson:"departure"`
	Arrival       string `json:"arrival" bson:"arrival"`
	Airline       string `json:"airline" bson:"airline"`
	DepartureTime string `json:"departure_time" bson:"departure_time"`
	ArrivalTime   string `json:"arrival_time" bson:"arrival_time"`
}

type UpdateFlightParams struct {
	DepartureTime string               `json:"departure_time,omitempty" bson:"departure_time,omitempty"`
	ArrivalTime   string               `json:"arrival_time,omitempty" bson:"arrival_time,omitempty"`
	Seats         []primitive.ObjectID `json:"seats,omitempty" bson:"seats,omitempty"`
}

func NewFlightFromParams(params CreateFlightParams) (*Flight, error) {
	return &Flight{
		Arrival:       params.Arrival,
		Departure:     params.Departure,
		Airline:       params.Airline,
		DepartureTime: params.DepartureTime,
		ArrivalTime:   params.ArrivalTime,
		Seats:         []primitive.ObjectID{},
	}, nil
}
