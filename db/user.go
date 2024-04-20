package db

import (
	"context"

	"github.com/fabrizioperria/goflight/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStorer interface {
	CreateRandomUser() error
	GetUserById(ctx context.Context, id string) (*types.User, error)
	GetUsers(ctx context.Context) ([]*types.User, error)
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

func (db *MongoDbUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	var cursor *mongo.Cursor
	cursor, err := db.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var results []*types.User
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (db *MongoDbUserStore) CreateRandomUser() error {
	_, err := db.collection.InsertOne(context.Background(), types.User{
		FirstName: "John",
		LastName:  "Doe",
	})
	return err
}
