package db

import (
	"context"

	"github.com/fabrizioperria/goflight/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	CreateRandomUser() error
	GetUserById(ctx context.Context, id string) (*types.User, error)
}

const (
	userCollection = "users"
)

type MongoDbUserStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDbUserStore(client *mongo.Client) *MongoDbUserStore {
	return &MongoDbUserStore{
		client:     client,
		collection: client.Database(dbName).Collection(userCollection),
	}
}

func (db *MongoDbUserStore) GetUserById(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user *types.User
	if err := db.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (db *MongoDbUserStore) CreateRandomUser() error {
	_, err := db.collection.InsertOne(context.Background(), types.User{
		FirstName: "John",
		LastName:  "Doe",
	})
	return err
}
