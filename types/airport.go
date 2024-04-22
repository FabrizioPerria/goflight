package types

type Airport struct {
	City string `json:"name" bson:"name"`
	Code string `json:"code" bson:"code"`
}
