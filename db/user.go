package db

import (
	"context"
	"fmt"

	"github.com/fabrizioperria/goflight/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Dropper interface {
	Drop(ctx context.Context) error
}

type UserStorer interface {
	CreateUser(ctx context.Context, user *types.User) (*types.User, error)
	GetUser(ctx context.Context, filter Map) (*types.User, error)
	GetUsers(ctx context.Context, pagination *Pagination) ([]*types.User, error)
	DeleteUser(ctx context.Context, filter Map) (string, error)
	UpdateUser(ctx context.Context, filter Map, values types.UpdateUserParams) (string, error)
	Dropper
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
		collection: client.Database(DBNAME).Collection(userCollection),
	}
}

func (db *MongoDbUserStore) GetUser(ctx context.Context, filter Map) (*types.User, error) {
	user := &types.User{}
	if err := db.collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (db *MongoDbUserStore) GetUsers(ctx context.Context, pagination *Pagination) ([]*types.User, error) {
	var cursor *mongo.Cursor
	cursor, err := db.collection.Find(ctx, Map{}, pagination.ToFindOptions())
	if err != nil {
		return nil, err
	}

	results := make([]*types.User, 0)
	err = cursor.All(ctx, &results)

	return results, err
}

func (db *MongoDbUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	result, err := db.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Id = result.InsertedID.(primitive.ObjectID)
	return user, err
}

func (db *MongoDbUserStore) DeleteUser(ctx context.Context, filter Map) (string, error) {
	res, err := db.collection.DeleteOne(ctx, filter)
	if err != nil || res.DeletedCount == 0 {
		return "", fmt.Errorf("user not found")
	}

	return "", nil
}

func (db *MongoDbUserStore) Drop(ctx context.Context) error {
	err := db.collection.Drop(ctx)
	return err
}

func (db *MongoDbUserStore) UpdateUser(ctx context.Context, filter Map, values types.UpdateUserParams) (string, error) {
	result, err := db.collection.UpdateOne(ctx, filter, Map{"$set": values})
	if err != nil || result.ModifiedCount == 0 {
		return "", fmt.Errorf("user not found")
	}
	return "", nil
}
