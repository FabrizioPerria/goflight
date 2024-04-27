package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type FlightTime struct {
	Day   int `json:"day" bson:"day"`
	Month int `json:"month" bson:"month"`
	Year  int `json:"year" bson:"year"`
	Hour  int `json:"hour" bson:"hour"`
	Min   int `json:"min" bson:"min"`
}

type Flight struct {
	Departure     string               `json:"departure" bson:"departure"`
	Arrival       string               `json:"arrival" bson:"arrival"`
	Airline       string               `json:"airline" bson:"airline"`
	Seats         []primitive.ObjectID `json:"seats" bson:"seats"`
	DepartureTime FlightTime           `json:"departure_time" bson:"departure_time"`
	ArrivalTime   FlightTime           `json:"arrival_time" bson:"arrival_time"`
	Id            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
}

type CreateFlightParams struct {
	Departure     string     `json:"departure" bson:"departure"`
	Arrival       string     `json:"arrival" bson:"arrival"`
	Airline       string     `json:"airline" bson:"airline"`
	DepartureTime FlightTime `json:"departure_time" bson:"departure_time"`
	ArrivalTime   FlightTime `json:"arrival_time" bson:"arrival_time"`
}

type UpdateFlightParams struct {
	Seats         []primitive.ObjectID `json:"seats,omitempty" bson:"seats,omitempty"`
	DepartureTime FlightTime           `json:"departure_time,omitempty" bson:"departure_time,omitempty"`
	ArrivalTime   FlightTime           `json:"arrival_time,omitempty" bson:"arrival_time,omitempty"`
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
