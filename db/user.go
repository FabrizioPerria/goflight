package db

import (
	"context"
	"fmt"

	"github.com/fabrizioperria/goflight/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStorer interface {
	CreateUser(ctx context.Context, user *types.User) (*types.User, error)
	GetUserById(ctx context.Context, id string) (*types.User, error)
	GetUsers(ctx context.Context) ([]*types.User, error)
	DeleteUserById(ctx context.Context, id string) (string, error)
	DeleteUsers(ctx context.Context) error
	// UpdateUser(ctx context.Context) error
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
	user := &types.User{}
	err = db.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user)

	return user, err
}

func (db *MongoDbUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	var cursor *mongo.Cursor
	cursor, err := db.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	results := make([]*types.User, 0)
	err = cursor.All(ctx, &results)

	return results, err
}

func (db *MongoDbUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	result, err := db.collection.InsertOne(ctx, user)
	user.Id = result.InsertedID.(primitive.ObjectID).Hex()
	return user, err
}

func (db *MongoDbUserStore) DeleteUserById(ctx context.Context, id string) (string, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}
	res, err := db.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil || res.DeletedCount == 0 {
		return "", fmt.Errorf("user %s not found", id)
	}

	return id, nil
}

func (db *MongoDbUserStore) DeleteUsers(ctx context.Context) error {
	_, err := db.collection.DeleteMany(ctx, bson.M{})
	return err
}
