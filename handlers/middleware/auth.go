package middleware

import (
	"os"
	"time"

	"github.com/fabrizioperria/goflight/types"
	"github.com/golang-jwt/jwt/v5"
)

func ProduceToken(user *types.User) string {
	now := time.Now()
	exp := now.Add(time.Hour * 4).UTC().Unix()
	claims := jwt.MapClaims{
		"sub": user.Id.Hex(),
		"exp": exp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	return signedToken
}
