package types

import (
	"golang.org/x/text/currency"
)

type SeatClass int

const (
	_ SeatClass = iota
	Economy
	Business
	First
)

type FlightTime struct {
	Day   int `json:"day" bson:"day"`
	Month int `json:"month" bson:"month"`
	Year  int `json:"year" bson:"year"`
	Hour  int `json:"hour" bson:"hour"`
	Min   int `json:"min" bson:"min"`
}

type Seat struct {
	Price     currency.Amount `json:"price" bson:"price"`
	Id        string          `json:"id,omitempty" bson:"_id,omitempty"`
	Number    string          `json:"number" bson:"number"`
	FlightId  string          `json:"flight_id" bson:"flight_id"`
	Class     SeatClass       `json:"class" bson:"class"`
	Available bool            `json:"available" bson:"available"`
}

type Airport struct {
	City string `json:"name" bson:"name"`
	Code string `json:"code" bson:"code"`
}

type Flight struct {
	Departure     Airport    `json:"departure" bson:"departure"`
	Arrival       Airport    `json:"arrival" bson:"arrival"`
	Id            string     `json:"id,omitempty" bson:"_id,omitempty"`
	Airline       string     `json:"airline" bson:"airline"`
	Seats         []Seat     `json:"seats" bson:"seats"`
	DepartureTime FlightTime `json:"departure_time" bson:"departure_time"`
	ArrivalTime   FlightTime `json:"arrival_time" bson:"arrival_time"`
}

type CreateFlightParams struct {
	Departure     Airport    `json:"departure" bson:"departure"`
	Arrival       Airport    `json:"arrival" bson:"arrival"`
	Airline       string     `json:"airline" bson:"airline"`
	DepartureTime FlightTime `json:"departure_time" bson:"departure_time"`
	ArrivalTime   FlightTime `json:"arrival_time" bson:"arrival_time"`
}

func NewFlightFromParams(params CreateFlightParams) Flight {
	return Flight{
		Arrival:       params.Arrival,
		Departure:     params.Departure,
		Airline:       params.Airline,
		DepartureTime: params.DepartureTime,
		ArrivalTime:   params.ArrivalTime,
	}
}
