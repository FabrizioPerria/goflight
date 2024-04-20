package types

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

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

const (
	bcryptCost         = 10
	minFirstNameLength = 3
	minLastNameLength  = 3
	minPasswordLength  = 8
)

func isValidEmail(email string) bool {
	validEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return validEmail.MatchString(email)
}

func (params CreateUserParams) Validate() []string {
	var errors []string
	if len(params.FirstName) < minFirstNameLength {
		errors = append(errors, fmt.Sprintf("first name must be at least %d characters", minFirstNameLength))
	}
	if len(params.LastName) < minLastNameLength {
		errors = append(errors, fmt.Sprintf("last name must be at least %d characters", minLastNameLength))
	}
	if len(params.PlainPassword) < minPasswordLength {
		errors = append(errors, fmt.Sprintf("password must be at least %d characters", minPasswordLength))
	}
	if !isValidEmail(params.Email) {
		errors = append(errors, "invalid email")
	}

	return errors
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
