package db

type Store struct {
	User        UserStorer
	Flight      FlightStorer
	Seat        SeatStorer
	Reservation ReservationStorer
}

func NewStore(user UserStorer, flight FlightStorer, seat SeatStorer, reservation ReservationStorer) *Store {
	return &Store{
		User:        user,
		Flight:      flight,
		Seat:        seat,
		Reservation: reservation,
	}
}
