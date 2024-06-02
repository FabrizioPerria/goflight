package db

import (
	"strconv"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Pagination struct {
	Page  string `json:"page"`
	Limit string `json:"limit"`
}

func (pagination *Pagination) ToFindOptions() *options.FindOptions {
	pageNumber, err := strconv.ParseInt(pagination.Page, 10, 64)
	if err != nil {
		pageNumber = 1
	}
	limitNumber, err := strconv.ParseInt(pagination.Limit, 10, 64)
	if err != nil {
		limitNumber = 10
	}
	opts := options.Find()
	opts.SetLimit(limitNumber)
	opts.SetSkip((pageNumber - 1) * limitNumber)
	return opts
}
