package types

type User struct {
	Id        string `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
}
