package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/fabrizioperria/goflight/types"
)

func SeedUsers() {
	fmt.Println("Seeding users")
	for i := 0; i < 1000; i++ {
		user := types.CreateUserParams{
			Email:         gofakeit.Email(),
			PlainPassword: gofakeit.Password(true, true, true, true, false, 10),
			Phone:         gofakeit.Phone(),
			FirstName:     gofakeit.FirstName(),
			LastName:      gofakeit.LastName(),
		}
		userMarshal, err := json.Marshal(user)
		if err != nil {
			fmt.Println(err)
			continue
		}
		request, _ := http.NewRequest("POST", "http://localhost:5001/api/v1/user", bytes.NewReader(userMarshal))

		request.Header.Set("Content-Type", "application/json")
		_, err = http.DefaultClient.Do(request)
		if err != nil {
			continue
		}
	}
}

func SeedFlights() {
	fmt.Println("Seeding flights")
	for i := 0; i < 100; i++ {
		flight := types.CreateFlightParams{
			Airline: gofakeit.Company(),
			Departure: types.Airport{
				City: gofakeit.City(),
				Code: gofakeit.DigitN(4),
			},
			Arrival: types.Airport{
				City: gofakeit.City(),
				Code: gofakeit.DigitN(4),
			},
			DepartureTime: types.FlightTime{
				Day:   gofakeit.Day(),
				Month: gofakeit.Month(),
				Year:  gofakeit.Year(),
				Hour:  gofakeit.Hour(),
				Min:   gofakeit.Minute(),
			},
			ArrivalTime: types.FlightTime{
				Day:   gofakeit.Day(),
				Month: gofakeit.Month(),
				Year:  gofakeit.Year(),
				Hour:  gofakeit.Hour(),
				Min:   gofakeit.Minute(),
			},
		}
		flightMarshal, err := json.Marshal(flight)
		if err != nil {
			fmt.Println(err)
			continue
		}
		request, _ := http.NewRequest("POST", "http://localhost:5001/api/v1/flight", bytes.NewReader(flightMarshal))

		request.Header.Set("Content-Type", "application/json")
		_, err = http.DefaultClient.Do(request)
		if err != nil {
			continue
		}
	}
}

func main() {
	SeedUsers()
	SeedFlights()
}
