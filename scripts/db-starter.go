package scripts

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
