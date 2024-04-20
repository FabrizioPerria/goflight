package types

import "golang.org/x/crypto/bcrypt"

type CreateUserParams struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	PlainPassword string `json:"password"`
}

type User struct {
	Id                string `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName         string `json:"first_name" bson:"first_name"`
	LastName          string `json:"last_name" bson:"last_name"`
	Email             string `json:"email" bson:"email"`
	Phone             string `json:"phone" bson:"phone"`
	EncryptedPassword string `json:"-" bson:"encrypted_password"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encrypted_password, err := bcrypt.GenerateFromPassword([]byte(params.PlainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		Phone:             params.Phone,
		EncryptedPassword: string(encrypted_password),
	}, nil
}
