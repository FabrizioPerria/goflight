package types

import (
	"time"

	"golang.org/x/text/currency"
)

type SeatClass int

const (
	_ SeatClass = iota
	Economy
	Business
	First
)

type Seat struct {
	Price     currency.Amount `json:"price" bson:"price"`
	Id        string          `json:"id,omitempty" bson:"_id,omitempty"`
	Number    string          `json:"number" bson:"number"`
	FlightId  string          `json:"flight_id" bson:"flight_id"`
	Class     SeatClass       `json:"class" bson:"class"`
	Available bool            `json:"available" bson:"available"`
}

type Airport struct {
	Id   string `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Code string `json:"code" bson:"code"`
}

type Flight struct {
	DepartureTime time.Time `json:"departure_time" bson:"departure_time"`
	ArrivalTime   time.Time `json:"arrival_time" bson:"arrival_time"`
	Departure     Airport   `json:"departure" bson:"departure"`
	Arrival       Airport   `json:"arrival" bson:"arrival"`
	Id            string    `json:"id,omitempty" bson:"_id,omitempty"`
	Airline       string    `json:"airline" bson:"airline"`
	Seats         []Seat    `json:"seats" bson:"seats"`
}
