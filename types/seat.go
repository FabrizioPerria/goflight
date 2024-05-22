package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type SeatClass int

const (
	_ SeatClass = iota
	Economy
	Business
	First
)

type SeatLocation int

const (
	_ SeatLocation = iota
	Aisle
	Middle
	window
)

type Seat struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FlightId  primitive.ObjectID `json:"flight_id" bson:"flight_id"`
	Number    int                `json:"number" bson:"number"`
	Price     float64            `json:"price" bson:"price"`
	Class     SeatClass          `json:"class" bson:"class"`
	Location  SeatLocation       `json:"location" bson:"location"`
	Available bool               `json:"available" bson:"available"`
}

type UpdateSeatParams struct {
	Price     float64 `json:"price" bson:"price"`
	Available bool    `json:"available" bson:"available"`
}
