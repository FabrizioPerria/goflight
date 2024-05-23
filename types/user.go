package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CreateUserParams struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	PlainPassword string `json:"password"`
}

type User struct {
	FirstName         string             `json:"first_name" bson:"first_name"`
	LastName          string             `json:"last_name" bson:"last_name"`
	Email             string             `json:"email" bson:"email"`
	Phone             string             `json:"phone" bson:"phone"`
	EncryptedPassword string             `json:"-" bson:"encrypted_password"`
	Id                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	IsAdmin           bool               `json:"is_admin" bson:"is_admin"`
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

func (user *User) Authenticate(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password))
	return err == nil
}

func (params CreateUserParams) Validate() map[string]string {
	errors := make(map[string]string)
	if len(params.FirstName) < minFirstNameLength {
		errors["first_name"] = fmt.Sprintf("first name must be at least %d characters", minFirstNameLength)
	}
	if len(params.LastName) < minLastNameLength {
		errors["last_name"] = fmt.Sprintf("last name must be at least %d characters", minLastNameLength)
	}
	if len(params.PlainPassword) < minPasswordLength {
		errors["password"] = fmt.Sprintf("password must be at least %d characters", minPasswordLength)
	}
	if !isValidEmail(params.Email) {
		errors["email"] = fmt.Sprintf("email %s is not valid", params.Email)
	}

	return errors
}

func (params UpdateUserParams) Validate() map[string]string {
	errors := make(map[string]string)

	if len(params.FirstName) < minFirstNameLength {
		errors["first_name"] = fmt.Sprintf("first name must be at least %d characters", minFirstNameLength)
	}
	if len(params.LastName) < minLastNameLength {
		errors["last_name"] = fmt.Sprintf("last name must be at least %d characters", minLastNameLength)
	}
	return errors
}

func NewUserFromParams(params CreateUserParams, isAdmin bool) (*User, error) {
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
		IsAdmin:           isAdmin,
	}, nil
}
