package types

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
	Id        string    `json:"id,omitempty" bson:"_id,omitempty"`
	FlightId  string    `json:"flight_id" bson:"flight_id"`
	Number    int       `json:"number" bson:"number"`
	Price     float64   `json:"price" bson:"price"`
	Class     SeatClass `json:"class" bson:"class"`
	Available bool      `json:"available" bson:"available"`
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

type UpdateFlightParams struct {
	DepartureTime FlightTime `json:"departure_time,omitempty" bson:"departure_time,omitempty"`
	ArrivalTime   FlightTime `json:"arrival_time,omitempty" bson:"arrival_time,omitempty"`
	Seats         []Seat     `json:"seats,omitempty" bson:"seats,omitempty"`
}

func NewFlightFromParams(params CreateFlightParams) (*Flight, error) {
	return &Flight{
		Arrival:       params.Arrival,
		Departure:     params.Departure,
		Airline:       params.Airline,
		DepartureTime: params.DepartureTime,
		ArrivalTime:   params.ArrivalTime,
	}, nil
}
