package db

import (
	"strconv"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Pagination struct {
	Page  string `json:"page"`
	Limit string `json:"limit"`
}

func (pagination *Pagination) GetPage() int64 {
	pageNumber, err := strconv.ParseInt(pagination.Page, 10, 64)
	if err != nil {
		return 1
	}
	return pageNumber
}

func (pagination *Pagination) GetLimit() int64 {
	limitNumber, err := strconv.ParseInt(pagination.Limit, 10, 64)
	if err != nil {
		return 10
	}
	return limitNumber
}

func (pagination *Pagination) ToFindOptions() *options.FindOptions {
	pageNumber := pagination.GetPage()
	limitNumber := pagination.GetLimit()
	opts := options.Find()
	opts.SetLimit(limitNumber)
	opts.SetSkip((pageNumber - 1) * limitNumber)
	return opts
}
