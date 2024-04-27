package db

type Store struct {
	User   UserStorer
	Flight FlightStorer
	Seat   SeatStorer
}

func NewStore(user UserStorer, flight FlightStorer, seat SeatStorer) *Store {
	return &Store{
		User:   user,
		Flight: flight,
		Seat:   seat,
	}
}
